import type { Config } from 'drizzle-kit';

export default {
  schema: './src/lib/server/schema.ts',
  driver: 'pg',
  dbCredentials: {
    connectionString: process.env.DB_URL,
  },
  schemaFilter: ['email', 'profile'],
  out: './src/lib/server/drizzle'
} satisfies Config;
