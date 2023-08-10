export interface Scope {
  id: string;
  description: string;
}

export interface ResourceServer {
  title: string;
  id: string;
  src: string;
  scopes: Scope[];
}
