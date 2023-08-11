import type { MyCredentialCreationOptions, MyCredentialRequestOptions } from './secure-client';
import * as secureClient from './secure-client';

export interface SuccessfulResponse {
  accessToken: string;
  refreshToken: string;
}

export const bufferToBase64 = (buffer: ArrayBuffer): string => btoa(String.fromCharCode(...new Uint8Array(buffer)));
export const base64ToBuffer = (base64: string): ArrayBuffer => Uint8Array.from(atob(base64), (c) => c.charCodeAt(0)).buffer;

export const register = async (credentialOptions: MyCredentialCreationOptions): Promise<void> => {
  const credential = await createCredential(credentialOptions);
  await secureClient.register(credentialOptions.authenticationId, credential);
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

export const login = async (credentialOptions: MyCredentialRequestOptions): Promise<void> => {
  const credential = await getCredential(credentialOptions);
  await secureClient.login(credentialOptions.authenticationId, credential);
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