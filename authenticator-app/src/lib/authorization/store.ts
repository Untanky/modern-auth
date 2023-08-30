import { browser } from '$app/environment';
import { refreshToken } from '$lib/secure-client';
import { readonly, writable } from 'svelte/store';
import {
    getAuthorizationState, setAuthorizationState,
    type AuthorizationData,
} from './local-store';

class AuthorizationState implements AuthorizationData {
    accessToken: string;
    refreshToken?: string;
    userId: string;
    expiresAt: number;

    constructor({
        accessToken, refreshToken, userId, expiresAt,
    }: AuthorizationData) {
        this.accessToken = accessToken;
        this.refreshToken = refreshToken;
        this.userId = userId;
        this.expiresAt = expiresAt;
    }

    refresh(): Promise<AuthorizationState | null> {
        if (!this.refreshToken) {
            return Promise.resolve(null);
        }

        return refreshToken(this.refreshToken)
            .then(({
                access_token, refresh_token, expires_in,
            }): AuthorizationState => new AuthorizationState({
                accessToken: access_token,
                refreshToken: refresh_token,
                expiresAt: Date.now() + expires_in,
                userId: this.userId,
            }))
            .catch(() => null);
    }
}

const internalStore = writable<AuthorizationState | null>();

export const initializeStoreLocally = () => {
    if (!browser) {
        return;
    }

    const state = getAuthorizationState();
    const newState = state ? new AuthorizationState(state) : null;

    internalStore.set(newState);

    internalStore.subscribe((state) => {
        setAuthorizationState(state);
    });

    internalStore.subscribe((state) => {
        if (state && state.refreshToken) {
            registerRefreshTimeout(state);
        }
    });
};

const FRESH_UNTIL = 60_000;

const registerRefreshTimeout = (state: AuthorizationState): NodeJS.Timeout => {
    const timeRemaing = state.expiresAt - Date.now() - FRESH_UNTIL;
    return setTimeout(() => state.refresh().then((state) => internalStore.set(state)), timeRemaing);
};

export const afterAuthentication = (data: AuthorizationData): void => {
    internalStore.set(new AuthorizationState(data));
};

export const afterLogout = (): void => {
    internalStore.set(null);
};

export const authorizationStore = readonly<AuthorizationData | null>(internalStore);
