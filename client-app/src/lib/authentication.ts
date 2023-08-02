export type MyCredentialCreationOptions = CredentialCreationOptions & { authenticationId: string; type: 'create' };
export type MyCredentialRequestOptions = CredentialRequestOptions & { authenticationId: string; type: 'get' };
export type CredentialOptions = MyCredentialCreationOptions | MyCredentialRequestOptions;

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
  await postNewCredential(credentialOptions.authenticationId, credential);
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

const postNewCredential = (authenticationId: string, credential: PublicKeyCredential): Promise<Response> => {
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
  });
};

export const signIn = async (credentialOptions: MyCredentialRequestOptions): Promise<void> => {
  const credential = await getCredential(credentialOptions);
  console.log(credential);
  validateCredential(credentialOptions.authenticationId, credential);
};

const getCredential = (credOps: CredentialRequestOptions): Promise<PublicKeyCredential> => {
  return navigator.credentials.get({
    publicKey: {
      ...credOps.publicKeyFoo,
      challenge: base64ToBuffer(credOps.publicKeyFoo.challenge as unknown as string),
      allowCredentials: credOps.publicKeyFoo.allowCredentials?.map((cred) => ({
        ...cred,
        id: base64ToBuffer(cred.id as unknown as string),
        transports: [],
      })),
    },
  }) as Promise<PublicKeyCredential>;
};

const validateCredential = (authenticationId: string, credential: PublicKeyCredential): Promise<Response> => {
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
  })
};
