import type { Preferences } from '$lib/preferences/model';

type verificationTemplate = {
    readonly type: 'verification';
    readonly props?: {
        doiLink: string;
        userName?: string;
        firstName?: string;
        lastName?: string;
    };
}

type accountResetTemplate = {
    readonly type: 'accountReset';
    readonly props?: {
        resetLink: string;
        userName?: string;
        firstName?: string;
        lastName?: string;
    };
}

type sessionNotificationTemplate = {
    readonly type: 'sessionNotification';
    readonly props?: {
        loginTime: string;
        loginLocation: string;
        userName?: string;
        firstName?: string;
        lastName?: string;
    };
}

type TemplateWithOptionalProps = verificationTemplate | accountResetTemplate | sessionNotificationTemplate;
export type Template = Required<TemplateWithOptionalProps>;
export type TemplateTypes = Template['type'];
export type TemplateProps = Template['props'];

export interface Email {
    readonly id: string;
    readonly sub: string;
    readonly template: TemplateWithOptionalProps;
    sentAt?: Date;
    deliveryMethod: 'resend';
    resendId: string;
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
