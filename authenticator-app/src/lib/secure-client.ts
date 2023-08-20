import { bufferToBase64 } from './authentication';

export type MyCredentialCreationOptions = CredentialCreationOptions & { authenticationId: string; type: 'create' };
export type MyCredentialRequestOptions = CredentialRequestOptions & { authenticationId: string; type: 'get' };
export type CredentialOptions = MyCredentialCreationOptions | MyCredentialRequestOptions;

export interface Grant {
  scope: string,
  client_id: string,
  sub: string,
  iat: number,
  exp: number,
  nbf?: number,
}

interface TokenResponse {
	access_token: string,
	token_type: string,
	expires_in: number,
	scope: string,
	refresh_token?: string
}

let correlationId: string;

export const initiateAuthentication = (userId: string): Promise<CredentialOptions> => {
    correlationId = crypto.randomUUID();

    return fetch('/v1/webauthn/authentication/initiate', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Request-Id': crypto.randomUUID(),
            'Correlation-Id': correlationId,
            'Cache-Control': 'no-store',
        },
        body: JSON.stringify({ userId }),
    }).then((response) => response.json());
};

export const register = (authenticationId: string, credential: PublicKeyCredential): Promise<TokenResponse> => {
    const { clientDataJSON, attestationObject } = credential.response as AuthenticatorAttestationResponse;

    // TODO: rename to /v1/authentication/register
    return fetch('/v1/webauthn/authentication/create', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Request-Id': crypto.randomUUID(),
            'Correlation-Id': correlationId,
            'Cache-Control': 'no-store',
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

export const login = (authenticationId: string, credential: PublicKeyCredential): Promise<TokenResponse> => {
    const {
        clientDataJSON, authenticatorData, signature, userHandle,
    } = credential.response as AuthenticatorAssertionResponse;

    // TODO: rename to /v1/authentication/login
    return fetch('/v1/webauthn/authentication/validate', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Request-Id': crypto.randomUUID(),
            'Correlation-Id': correlationId,
            'Cache-Control': 'no-store',
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
            },
        }),
    }).then((response) => response.json());
};

export const validate = (accessToken: string): Promise<Grant> => {
    return fetch('/v1/oauth/token/validate', {
        method: 'POST',
        headers: {
            'Authorization': `Bearer ${accessToken}`,
            'Content-Type': 'application/json',
            'Request-Id': crypto.randomUUID(),
            'Cache-Control': 'no-store',
        },
    }).then((response) => response.json());
};

export const refreshToken = (refreshToken: string): Promise<TokenResponse> => {
    const data = new URLSearchParams();
    data.append('refresh_token', refreshToken);
    data.append('grant_type', 'refresh_token');
    data.append('client_id', 'abc');

    return fetch(`/v1/oauth/token`, {
        method: 'POST',
        headers: {
            'Request-Id': crypto.randomUUID(),
            'Cache-Control': 'no-store',
        },
        body: data,
    }).then((res) => res.json());
};
