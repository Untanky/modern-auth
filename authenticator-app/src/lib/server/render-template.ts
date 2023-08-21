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

type TemplatePropsOf<Type extends TemplateTypes> = Extract<Template, { type: Type }>['props'];
type ComponentOf<Type extends TemplateTypes> = ComponentType<SvelteComponent<TemplatePropsOf<Type>>>

type TemplateToComponentsMapping = {
    [T in TemplateTypes]: ComponentOf<T>;
};

type TemplateToSubjectMapping = {
    [T in TemplateTypes]: (props?: TemplatePropsOf<T>) => string;
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
    const template = templateToComponents[type] as ComponentOf<typeof type>;
    const body = render({ template, props });
    const subject = templateToSubjectMapping[type]();

    return {
        body,
        subject,
    };
};
