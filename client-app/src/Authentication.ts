export const initiateAuthentication = async (userId: string): Promise<CredentialCreationOptions> => {
  return fetch('http://localhost:3000/v1/webauthn/authentication/initiate', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ userId }),
  }).then((response) => response.json());
};

export const createCredential = async () => {

};

export const getCredential = async () => [

];
