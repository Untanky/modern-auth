import { writable } from "svelte/store";

export type AuthorizationState = 'userId' | 'createCredential' | 'getCredential' | 'success';
  
let authorizationState: AuthorizationState = 'userId';
let loading = false;
let error: Error | null = null;
let userId: string | null = null;

type authorizationStateType = {
  state: 'userId';
  loading: boolean;
  error: Error | null;
} | {
  state: 'createCredential';
  loading: boolean;
  error: Error | null;
  userId: string;
  credentialOptions: CredentialCreationOptions;
} | {
  state: 'getCredential';
  loading: boolean;
  error: Error | null;
  userId: string;
  credentialOptions: CredentialRequestOptions;
} | {
  state: 'success';
  loading: boolean;
  error: Error | null;
  userId: string;
};

export const state = writable<authorizationStateType>({
  state: authorizationState,
  loading,
  error,
});
