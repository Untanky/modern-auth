export type MyCredentialCreationOptions = CredentialCreationOptions & { authenticationId: string; type: 'create' };
export type MyCredentialRequestOptions = CredentialRequestOptions & { authenticationId: string; type: 'get' };
export type CredentialOptions = MyCredentialCreationOptions | MyCredentialRequestOptions;

export interface SuccessfulResponse {
  accessToken: string;
  refreshToken: string;
}

export const initiateAuthentication = async (userId: string): Promise<CredentialOptions> => {
  return fetch('/v1/webauthn/authentication/initiate', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ userId }),
  }).then((response) => response.json());
};

const bufferToBase64 = (buffer: ArrayBuffer): string => btoa(String.fromCharCode(...new Uint8Array(buffer)));
const base64ToBuffer = (base64: string): ArrayBuffer => Uint8Array.from(atob(base64), (c) => c.charCodeAt(0)).buffer;

export const signUp = async (credentialOptions: MyCredentialCreationOptions): Promise<void> => {
  const credential = await createCredential(credentialOptions);
  const tokens = await postNewCredential(credentialOptions.authenticationId, credential);
  storeTokens(tokens);
};

const createCredential = (credOps: CredentialCreationOptions): Promise<PublicKeyCredential> => {
  return navigator.credentials.create({
    publicKey: {
      ...credOps.publicKey,
      challenge: base64ToBuffer(credOps.publicKey.challenge as unknown as string),
      user: {
        ...credOps.publicKey.user,
        id: base64ToBuffer(credOps.publicKey.user.id as unknown as string),
      },
    },
  }) as Promise<PublicKeyCredential>;
};

const postNewCredential = (authenticationId: string, credential: PublicKeyCredential): Promise<SuccessfulResponse> => {
  const { clientDataJSON, attestationObject } = credential.response as AuthenticatorAttestationResponse;

  return fetch('/v1/webauthn/authentication/create', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      id: credential.id,
      authenticationId,
      type: credential.type,
      response: {
        clientDataJSON: bufferToBase64(clientDataJSON),
        attestationObject: bufferToBase64(attestationObject),
      },
    }),
  }).then((response) => response.json());
};

export const signIn = async (credentialOptions: MyCredentialRequestOptions): Promise<void> => {
  const credential = await getCredential(credentialOptions);
  const tokens = await validateCredential(credentialOptions.authenticationId, credential);
  storeTokens(tokens);
};

const getCredential = (credOps: CredentialRequestOptions): Promise<PublicKeyCredential> => {
  return navigator.credentials.get({
    publicKey: {
      ...credOps.publicKey,
      challenge: base64ToBuffer(credOps.publicKey.challenge as unknown as string),
      allowCredentials: credOps.publicKey.allowCredentials?.map((cred) => ({
        ...cred,
        id: base64ToBuffer(cred.id as unknown as string),
        transports: [],
      })),
    },
  }) as Promise<PublicKeyCredential>;
};

const validateCredential = (authenticationId: string, credential: PublicKeyCredential): Promise<SuccessfulResponse> => {
  const { clientDataJSON, authenticatorData, signature, userHandle } = credential.response as AuthenticatorAssertionResponse;

  return fetch('/v1/webauthn/authentication/validate', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      id: credential.id,
      rawId: bufferToBase64(credential.rawId),
      authenticationId,
      type: credential.type,
      response: {
        clientDataJSON: bufferToBase64(clientDataJSON),
        authenticatorData: bufferToBase64(authenticatorData),
        signature: bufferToBase64(signature),
        userHandle: bufferToBase64(userHandle),
      }
    }),
  }).then((response) => response.json());
};

const getTokens = (): SuccessfulResponse => {
  const accessToken = localStorage.getItem('accessToken');
  const refreshToken = localStorage.getItem('refreshToken');
  
  return { accessToken, refreshToken };
};

const storeTokens = (tokens: SuccessfulResponse): void => {
  localStorage.setItem('accessToken', tokens.accessToken);
  localStorage.setItem('refreshToken', tokens.refreshToken);
};
