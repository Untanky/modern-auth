import { readonly, writable } from "svelte/store";

interface AuthenticationData {
  accessToken: string;
  refreshToken?: string;
  userId: string;
}

const getAuthenticationStateOrNull = (): AuthenticationData | null => {
  const accessToken = localStorage.getItem('accessToken');
  const refreshToken = localStorage.getItem('refreshToken');
  const userId = localStorage.getItem('userId');

  if (!accessToken || !userId || !refreshToken) {
    return null;
  }

  return {
    accessToken,
    refreshToken,
    userId,
  }
}

const initialData = getAuthenticationStateOrNull();

const state = writable<AuthenticationData | null>(initialData, () => {

});

export const onAuthenticated = (data: AuthenticationData): void => {
  localStorage.set('accessToken', data.accessToken);
  localStorage.set('userId', data.userId);
  if (data.refreshToken) {
    localStorage.set('refreshToken', data.refreshToken);
  }

  state.set(data);
};

export const onLoggedOut = (): void => {
  localStorage.removeItem('accessToken');
  localStorage.removeItem('refreshToken');

  state.set(null);
}

export const authenticationState = readonly(state);
