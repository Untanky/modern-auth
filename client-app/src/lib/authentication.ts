export const initiateAuthentication = async (userId: string): Promise<CredentialCreationOptions> => {
  return fetch('/v1/webauthn/authentication/initiate', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ userId }),
  }).then((response) => response.json());
};

const utf8Decoder = new TextDecoder('utf-8');
const createUint8ArrayFrom = (value: string): Uint8Array => Uint8Array.from(value, (c) => c.charCodeAt(0));

export const signUp = async (credentialOptions: CredentialCreationOptions): Promise<void> => {
  const credential = await createCredential(credentialOptions);
  await postNewCredential(credential);
};

const createCredential = (credOps: CredentialCreationOptions): Promise<PublicKeyCredential> => {
  return navigator.credentials.create({
    publicKey: {
      ...credOps.publicKey,
      challenge: createUint8ArrayFrom(credOps.publicKey.challenge as unknown as string),
      user: {
        ...credOps.publicKey.user,
        id: createUint8ArrayFrom(credOps.publicKey.user.id as unknown as string),
      },
    },
  }) as Promise<PublicKeyCredential>;
};

const postNewCredential = (credential: PublicKeyCredential): Promise<Response> => {
  const { clientDataJSON, attestationObject } = credential.response as AuthenticatorAttestationResponse;

  return fetch('/v1/webauthn/authentication/create', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      id: credential.id,
      type: credential.type,
      response: {
        clientDataJSON: utf8Decoder.decode(clientDataJSON),
        attestationObject: utf8Decoder.decode(attestationObject),
      },
    }),
  });
};

export const signIn = async (credentialOptions: CredentialRequestOptions): Promise<void> => {
  const credential = await getCredential(credentialOptions);
  validateCredential(credential);
};

const getCredential = (credOps: CredentialRequestOptions): Promise<PublicKeyCredential> => {
  return navigator.credentials.get({
    publicKey: {
      ...credOps.publicKey,
      challenge: createUint8ArrayFrom(credOps.publicKey.challenge as unknown as string),
      allowCredentials: credOps.publicKey.allowCredentials?.map((cred) => ({
        ...cred,
        id: createUint8ArrayFrom(cred.id as unknown as string),
      })),
    },
  }) as Promise<PublicKeyCredential>;
};

const validateCredential = (credential: PublicKeyCredential): Promise<Response> => {
  const { clientDataJSON, authenticatorData, signature, userHandle } = credential.response as AuthenticatorAssertionResponse;

  return fetch('/v1/webauthn/authentication/validate', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      id: credential.id,
      type: credential.type,
      response: {
        clientDataJSON: utf8Decoder.decode(clientDataJSON),
        authenticatorData: utf8Decoder.decode(authenticatorData),
        signature: utf8Decoder.decode(signature),
        userHandle: utf8Decoder.decode(userHandle),
      }
    }),
  })
};
