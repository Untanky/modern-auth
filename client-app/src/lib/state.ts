import { writable } from "svelte/store";
import type { MyCredentialCreationOptions, MyCredentialRequestOptions } from "./secure-client";

export type AuthorizationState = 'userId' | 'createCredential' | 'getCredential' | 'success';
  
let authorizationState: AuthorizationState = 'userId';
let loading = false;
let error: Error | null = null;

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
  state: authorizationState,
  loading,
  error,
});
