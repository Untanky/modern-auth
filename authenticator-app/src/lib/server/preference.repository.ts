import type { Preferences } from '$lib/preferences/model';
import { eq, sql } from 'drizzle-orm';
import type { PostgresJsDatabase } from 'drizzle-orm/postgres-js';
import type { PreferencesRepository } from './email.model';
import type * as schema from './schema';
import { preferences } from './schema';

export class DrizzlePreferencesRepository implements PreferencesRepository {
    private readonly db: PostgresJsDatabase<typeof schema>;

    constructor(db: PostgresJsDatabase<typeof schema>) {
        this.db = db;
    }

    findFirst(where?: Partial<Preferences> | undefined): Promise<Preferences> {
        return this.db
            .select()
            .from(preferences)
            .where(where && where.sub ? eq(preferences.sub, where.sub) : sql`1 = 1`)
            .limit(1)
            .then(([result]): Preferences => ({
                ...result,
                verifiedAt: result.verifiedAt || undefined,
            }));
    }

    findMany(): Promise<Preferences[]> {
        throw new Error('Method not implemented.');
    }

    create(entity: Preferences): Promise<Preferences> {
        return this.db
            .insert(preferences)
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
        return this.db.update(preferences)
            .set(entity)
            .where(eq(preferences.sub, entity.sub))
            .returning()
            .then(([result]) => ({
                ...result,
                // TODO: find a way to fetch these dynamically
                verified: false,
                verifiedAt: undefined,
            }));
    }

    delete(where: Partial<Preferences>): Promise<void> {
        return this.db.delete(preferences)
            .where(where && where.sub ? eq(preferences.sub, where.sub) : sql`1 = 1`)
            .then();
    }
}
