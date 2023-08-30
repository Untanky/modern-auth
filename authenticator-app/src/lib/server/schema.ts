/* eslint-disable newline-per-chained-call */
import {
    boolean, char, pgEnum, pgSchema, timestamp, uuid, varchar,
} from 'drizzle-orm/pg-core';
import { CODE_LENGTH } from './verification.service';

export const emailSchema = pgSchema('email');

export const preference = emailSchema.table('preference', {
    sub: uuid('sub').primaryKey()
        .defaultRandom(),
    emailAddress: varchar('email_address').notNull(),
    allowAccountReset: boolean('allow_account_reset').default(true)
        .notNull(),
    allowSessionNotification: boolean('allow_session_notification').default(false)
        .notNull(),
});

export const templates = pgEnum('templates_enum', [
    'verification',
    'accountReset',
    'sessionNotification',
]);

export const email = emailSchema.table('email', {
    id: uuid('id').primaryKey().defaultRandom(),
    sub: uuid('sub').references(() => preference.sub, { onDelete: 'no action' }).notNull(),
    sentAt: timestamp('sent_at').notNull(),
    template: templates('template').notNull(),
});

export const verificationRequest = emailSchema.table('verification_request', {
    id: uuid('id').primaryKey().defaultRandom(),
    sub: uuid('sub').references(() => preference.sub, { onDelete: 'no action' }).notNull(),
    expiresAt: timestamp('expires_at').notNull(),
    codeVerifier: char('code_verifier', { length: CODE_LENGTH }).notNull(),
});

export const verification = emailSchema.table('verification', {
    id: uuid('id').primaryKey().references(() => verificationRequest.id, { onDelete: 'cascade' }),
    verifiedAt: timestamp('verified_at').notNull(),
});

export const resendEmail = emailSchema.table('resend_email', {
    id: uuid('id').primaryKey().references(() => email.id, { onDelete: 'cascade' }),
    resendId: varchar('resend_id').notNull().unique(),
});
