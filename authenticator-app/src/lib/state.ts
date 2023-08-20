import { writable } from 'svelte/store';
import type { MyCredentialCreationOptions, MyCredentialRequestOptions } from './secure-client';

export type AuthorizationState = 'userId' | 'createCredential' | 'getCredential' | 'success';

type authorizationStateType = {
  state: 'userId';
  loading: boolean;
  error: Error | null;
} | {
  state: 'createCredential';
  loading: boolean;
  error: Error | null;
  userId: string;
  credentialOptions: MyCredentialCreationOptions;
} | {
  state: 'getCredential';
  loading: boolean;
  error: Error | null;
  userId: string;
  credentialOptions: MyCredentialRequestOptions;
} | {
  state: 'success';
  loading: boolean;
  error: Error | null;
  userId: string;
};

export const state = writable<authorizationStateType>({
    state: 'userId',
    loading: false,
    error: null,
});
