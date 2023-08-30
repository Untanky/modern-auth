import type { PostgresJsDatabase } from 'drizzle-orm/postgres-js';
import type { PreferencesRepository } from './email.model';
import type * as schema from './schema';
import {
    preference,
    verification,
    verificationRequest,
} from './schema';
import { eq, sql } from 'drizzle-orm';
import type { Preferences } from '$lib/preferences/model';

export class DrizzlePreferencesRepository implements PreferencesRepository {
    private readonly db: PostgresJsDatabase<typeof schema>;

    constructor(db: PostgresJsDatabase<typeof schema>) {
        this.db = db;
    }

    findFirst(where?: Partial<Preferences> | undefined): Promise<Preferences> {
        return this.db
            .select()
            .from(preference)
            .leftJoin(verificationRequest, eq(preference.sub, verificationRequest.sub))
            .leftJoin(verification, eq(verificationRequest.id, verification.id))
            .where(where && where.sub ? eq(preference.sub, where.sub) : sql`1 = 1`)
            .limit(1)
            .then(([result]): Preferences => ({
                sub: result.preference.sub,
                allowAccountReset: result.preference.allowAccountReset,
                allowSessionNotification: result.preference.allowSessionNotification,
                emailAddress: result.preference.emailAddress,
                verified: !!result.verification,
                verifiedAt: result.verification?.verifiedAt,
            }));
    }

    findMany(): Promise<Preferences[]> {
        throw new Error('Method not implemented.');
    }

    create(entity: Preferences): Promise<Preferences> {
        return this.db
            .insert(preference)
            .values(entity)
            .returning()
            .then(([result]) => ({
                ...result,
                // TODO: find a way to fetch these dynamically
                verified: false,
                verifiedAt: undefined,
            }));
    }

    update(entity: Preferences): Promise<Preferences> {
        return this.db.update(preference)
            .set(entity)
            .where(eq(preference.sub, entity.sub))
            .returning()
            .then(([result]) => ({
                ...result,
                // TODO: find a way to fetch these dynamically
                verified: false,
                verifiedAt: undefined,
            }));
    }

    delete(where: Partial<Preferences>): Promise<void> {
        return this.db.delete(preference)
            .where(where && where.sub ? eq(preference.sub, where.sub) : sql`1 = 1`)
            .then();
    }
}
