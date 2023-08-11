import { refreshToken } from "$lib/secure-client";
import { readonly, writable, type Updater, get } from "svelte/store";

interface AuthenticationData {
  accessToken: string;
  refreshToken?: string;
  userId: string;
}

const AUTH_KEY = 'auth_data';

const setAuthenticationData = (data: AuthenticationData | null): void => {
  if (!data) {
    localStorage.removeItem(AUTH_KEY);
    return
  }

  localStorage.setItem(AUTH_KEY, btoa(JSON.stringify(data)));
};

const getAuthenticationState = (): AuthentiationState => {
  const rawData = localStorage.getItem(AUTH_KEY);

  if (!rawData) {
    return Promise.resolve(null);
  }

  return Promise.resolve(JSON.parse(atob(rawData)));
}

const initialData = getAuthenticationState();

type AuthentiationState = Promise<AuthenticationData | null>;

const state = writable<AuthentiationState>(initialData, (_, update) => {
  const refresh: Updater<AuthentiationState> = (value) => {
    return value.then(data => {
      if (data && data.refreshToken) {
        return refreshToken(data?.refreshToken)
          .then(({ access_token, refresh_token }): AuthenticationData => ({
            accessToken: access_token,
            refreshToken: refresh_token,
            userId: data.userId,
          }))
          .catch(() => null);
      }
      return null;
    })
  };

  update(refresh);

  const intervalId = setInterval(() => {
    update(refresh);
  }, 60 * 60 * 1000, );

  return () => {
    clearInterval(intervalId);
  }
});

state.subscribe((data) => {
  data.then(setAuthenticationData);
});

export const onAuthenticated = (data: AuthenticationData): void => {
  state.set(Promise.resolve(data));
};

export const onLoggedOut = (): void => {
  state.set(Promise.resolve(null));
}

export const authenticationState = readonly(state);
