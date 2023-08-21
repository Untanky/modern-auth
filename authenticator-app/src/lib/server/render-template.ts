import { render } from 'svelte-email';
import Verification from './templates/Verification.svelte';

type TemplateResult = {
  body: string;
  subject: string;
};

export const renderTemplate = (): TemplateResult => {
    const body = render({ template: Verification, props: { name: 'Lukas', verificationLink: 'http://localhost:3000/' } });

    return {
        body,
        subject: 'TEST EMAIL',
    };
};
