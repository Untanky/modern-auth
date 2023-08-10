interface AuthenticationData {
  accessToken: string;
  refreshToken?: string;
  userId: string;
}

export const logout = () => {
  localStorage.removeItem('accessToken');
  localStorage.removeItem('refreshToken');
}

type Nullable<T> = { [P in keyof T]?: T[P] | null; }

export const getTokens = (): Nullable<AuthenticationData> => {
  const accessToken = localStorage.getItem('accessToken');
  const refreshToken = localStorage.getItem('refreshToken');
  const userId = localStorage.getItem('userId');
  
  return { accessToken, refreshToken, userId };
};

export const storeTokens = (tokens: AuthenticationData): void => {
  localStorage.setItem('accessToken', tokens.accessToken);
  localStorage.setItem('userId', tokens.userId);
  if (tokens.refreshToken) {
    localStorage.setItem('refreshToken', tokens.refreshToken);
  }
};
