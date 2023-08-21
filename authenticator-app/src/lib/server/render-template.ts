import type { ComponentType, SvelteComponent } from 'svelte';
import { render } from 'svelte-email';
import type { Template, TemplateTypes } from './email.model';
import Verification from './templates/Verification.svelte';
import AccountReset from './templates/AccountReset.svelte';
import SessionNotification from './templates/SessionNotification.svelte';

type TemplateResult = {
  body: string;
  subject: string;
};

type TemplateToComponentsMapping = {
    [T in TemplateTypes]: ComponentType<SvelteComponent<Extract<Template, { type: T }>['props']>>;
};

type TemplateToSubjectMapping = {
    [T in TemplateTypes]: (props?: Extract<Template, { type: T }>['props']) => string;
}

const templateToComponents = {
    verification: Verification,
    accountReset: AccountReset,
    sessionNotification: SessionNotification,
} satisfies TemplateToComponentsMapping;

const templateToSubjectMapping = {
    verification: () => 'Account Verification',
    accountReset: () => 'Reset Account',
    sessionNotification: () => 'New Login to Account',
} satisfies TemplateToSubjectMapping;

export const renderTemplate = ({ type, props }: Template): TemplateResult => {
    const template = templateToComponents[type] as ComponentType;
    const body = render({ template, props });
    const subject = templateToSubjectMapping[type]();

    return {
        body,
        subject,
    };
};
