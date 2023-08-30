export interface Preferences {
    readonly sub: string;
    emailAddress: string;
    verified: boolean;
    verifiedAt?: Date;
    allowAccountReset: boolean;
    allowSessionNotification: boolean;
}
