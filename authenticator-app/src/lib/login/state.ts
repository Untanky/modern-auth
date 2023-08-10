import { storeTokens } from "$lib/authorizationStore";
import { initiateAuthentication, login } from "$lib/secure-client";
import { readonly, writable } from "svelte/store"; 

abstract class AbstractState {
  loading = false
  error: Error | null = null;

  async identified(userId: string, options: CredentialRequestOptions): Promise<AbstractState> {
    throw new Error("Operation not supported");
  }

  async enteredPassword(password: string): Promise<AbstractState> {
    throw new Error("Operation not supported");
  }

  async enteredOtp(otp: string): Promise<AbstractState> {
    throw new Error("Operation not supported");
  }

  async validatedCredential(credential: PublicKeyCredential): Promise<AbstractState> {
    throw new Error("Operation not supported");
  }
}

class IdentificationState extends AbstractState {
  state = 'identification';

  async identified(userId: string): Promise<AbstractState> {
    this.loading = true;

    try {
      const options = await initiateAuthentication(userId);
      return new CredentialState(userId, options);
    } catch (e) {
      this.error = e as Error;
      return this;
    } finally {
      this.loading = false;
    }
  }
}

class CredentialState extends AbstractState {
  state = 'credential';

  constructor(public userId: string, public options: CredentialRequestOptions) {
    super();
  }

  async enteredPassword(password: string): Promise<AbstractState> {
    return new PasswordState(this.userId, this.options, password);
  }

  async validatedCredential(credential: PublicKeyCredential): Promise<AbstractState> {
    const tokenResponse = await login('', credential);
    storeTokens({
      accessToken: tokenResponse.access_token,
      refreshToken: tokenResponse.refresh_token,
      userId: this.userId,
    });
    
    return new SuccessState(this.userId);
  }
}

class PasswordState extends AbstractState {
  state = 'password';
  otp?: string;

  constructor(public userId: string, public options: CredentialRequestOptions, public password: string) {
    super();
  }

  async enteredOtp(otp: string): Promise<AbstractState> {
    this.otp = otp;
    return this;
  }

  async validatedCredential(credential: PublicKeyCredential): Promise<AbstractState> {
    const tokenResponse = await login('', credential);
    storeTokens({
      accessToken: tokenResponse.access_token,
      refreshToken: tokenResponse.refresh_token,
      userId: this.userId,
    });
    
    return new SuccessState(this.userId);
  }
}

class SuccessState extends AbstractState {
  state = 'success';

  constructor(public userId: string) {
    super();
  }
}

export type LoginState = IdentificationState | CredentialState | PasswordState | SuccessState;

const internalState = writable<LoginState>(new IdentificationState);

export const state = readonly(internalState);
