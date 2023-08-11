import { browser } from "$app/environment";

export interface AuthorizationData {
  accessToken: string;
  refreshToken?: string;
  userId: string;
  expiresAt: number;
}
type AuthorizationState = AuthorizationData | null;

const AUTH_KEY = 'auth_key';

export const getAuthorizationState = (): AuthorizationState => {
  if (!browser) {
    return null;
  }

  const rawData = localStorage.getItem(AUTH_KEY);
  if (!rawData) {
    return null;
  }

  return JSON.parse(atob(rawData));
};

export const setAuthorizationState = (state: AuthorizationState): void => {
  if (!browser) {
    throw new Error('Operation only supported on the user agent!');
  }

  if (!state) {
    localStorage.removeItem(AUTH_KEY);
  }

  localStorage.setItem(AUTH_KEY, btoa(JSON.stringify(localStorage)));
}