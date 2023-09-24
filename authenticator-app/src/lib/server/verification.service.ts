import dayjs from 'dayjs';
import type { InferSelectModel } from 'drizzle-orm';
import {
    and, eq, gte,
} from 'drizzle-orm';
import type { PostgresJsDatabase } from 'drizzle-orm/postgres-js';
// import { shake256 as hasher } from 'js-sha3';
import crypto from 'node:crypto';
import type { EmailService } from './email.service';
import type * as schema from './schema';
import { preferences, verificationRequest } from './schema';

export const CODE_LENGTH = 48;
const BITS_IN_NIBBLE = 4;

const BASE_URL = 'http://localhost:5173';

type VerificationParameters = {
    id: string;
    code: string;
}

type SelectVerificationRequest = InferSelectModel<typeof verificationRequest>;

const getVerificationLink = ({ id, code }: VerificationParameters): string => {
    return `${BASE_URL}/v1/verification/${id}?code=${code}`;
};

// const hash = (value: string): string => hasher(value, CODE_LENGTH * BITS_IN_NIBBLE);

const generateCodeVerifier = (): [string, string] => {
    const random = crypto.randomBytes(64).toString('hex');
    // const code = hash(random);
    // const codeVerifier = hash(code);
    return [
        '',
        '',
    ];
};

export class VerificationService {
    readonly #db: PostgresJsDatabase<typeof schema>;
    readonly #emailService: EmailService;

    constructor(db: PostgresJsDatabase<typeof schema>, emailService: EmailService) {
        this.#db = db;
        this.#emailService = emailService;
    }

    async startVerification(sub: string): Promise<void> {
        const [
            code,
            codeVerifier,
        ] = generateCodeVerifier();
        const { id } = await this.createRequest(sub, codeVerifier);
        await this.sendEmail(sub, { id, code });
    }

    private async createRequest(sub: string, codeVerifier: string): Promise<{ id: string }> {
        const expiresAt = dayjs().add(3, 'days');
        const [{ id }] = await this.#db
            .insert(verificationRequest)
            .values({
                sub,
                codeVerifier,
                expiresAt: expiresAt.toDate(),
            })
            .returning({ id: verificationRequest.id });
        return { id };
    }

    private async sendEmail(sub: string, params: VerificationParameters): Promise<void> {
        const link = getVerificationLink(params);
        await this.#emailService.send({
            sub,
            template: {
                type: 'verification',
                props: { doiLink: link },
            },
        });
    }

    async finishVerification(id: string, code: string): Promise<void> {
        const { sub } = await this.findValidVerificationRequest({ id, code });
        await this.updateVerificationState(sub);
    }

    private async findValidVerificationRequest({ id, code }: VerificationParameters): Promise<SelectVerificationRequest> {
        console.log(new Date());
        const [request] = await this.#db
            .select()
            .from(verificationRequest)
            .where(and(
                eq(verificationRequest.id, id),
                eq(verificationRequest.codeVerifier, hash(code)),
                gte(verificationRequest.expiresAt, new Date()),
            ))
            .limit(1);
        console.log(hash(code));
        return request;
    }

    private async updateVerificationState(sub: string): Promise<void> {
        await this.#db.update(preferences)
            .set({
                verified: true,
                verifiedAt: new Date(),
            })
            .where(eq(preferences.sub, sub));
    }
}
