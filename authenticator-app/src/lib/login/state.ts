import { readonly, writable } from 'svelte/store';

type IdentificationState = {
  state: 'identification';
}

type CredentialState = {
  state: 'credential';
  userId: string;
  options: CredentialRequestOptions;
}

type PasswordState = {
  state: 'password';
  userId: string;
  password: string;
  otp?: string;
}

type AuthorizationState = {
  state: 'authorization';
  userId: string;
  authorizationContext: any;
}

type SuccessState = {
  state: 'success';
  userId: string;
}

export type LoginState =
  | IdentificationState
  | CredentialState
  | PasswordState
  | AuthorizationState
  | SuccessState;

const internalState = writable<LoginState>({ state: 'identification' });

export const state = readonly(internalState);

type CreateCredentialState = {
  state: 'createCredential';
  userId: string;
  options: CredentialCreationOptions;
}

type CreatePasswordState = {
  state: 'createPassword';
  userId: string;
  options: CredentialCreationOptions;
}

type SetupOtpState = {
  state: 'setupOtp';
  userId: string;
  options: CredentialCreationOptions;
}

type SetupProfile = {
  state: 'setupProfile';
  userId: string;
}

type SetupEmail = {
  state: 'setupEmail';
  userId: string;
}

export type RegistrationState = 
  | IdentificationState 
  | CreateCredentialState 
  | CreatePasswordState 
  | SetupOtpState 
  | SetupProfile 
  | SetupEmail 
  | AuthorizationState
  | SuccessState;
