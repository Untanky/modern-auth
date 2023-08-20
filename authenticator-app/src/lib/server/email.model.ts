type verificationTemplate = {
  readonly type: 'verification';
  readonly sub: string;
  doiLink: string;
  userName?: string;
  firstName?: string;
  lastName?: string;
}

type accountResetTemplate = {
  readonly type: 'accountReset';
  readonly sub: string;
  resetLink: string;
  userName?: string;
  firstName?: string;
  lastName?: string;
}

type sessionNotificationTemplate = {
  readonly type: 'sessionNotification';
  readonly sub: string;
  loginTime: string;
  loginLocation: string;
  userName?: string;
  firstName?: string;
  lastName?: string;
}

export type Template = verificationTemplate | accountResetTemplate | sessionNotificationTemplate;
export type TemplateTypes = Template['type'];

export interface Email {
  readonly id: string;
  readonly sub: string;
  readonly template: Readonly<Template>;
  sentAt?: Date;
  deliveryMethod: 'resend';
  resendId: string;
}

export interface Preferences {
  readonly sub: string;
  email: string;
  verified: string;
  verifiedAt?: Date;
  allowAccountReset: boolean;
  allowSessionNotification: boolean;
}

interface Repository<Type> {
  findFirst(where?: Partial<Type>): Promise<Type>;
  findMany(where?: Partial<Type>): Promise<Type[]>;
  create(entity: Type): Promise<Type>;
  update(entity: Type): Promise<Type>;
  delete(where: Partial<Type>): Promise<void>;
}

export interface EmailRepository extends Repository<Email> {
  create(entity: Required<Email>): Promise<Email>;
}

export type PreferencesRepository = Repository<Preferences>;
