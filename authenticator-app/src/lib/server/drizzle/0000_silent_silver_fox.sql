CREATE SCHEMA "email";
--> statement-breakpoint
DO $$ BEGIN
 CREATE TYPE "templates_enum" AS ENUM('verification', 'accountReset', 'sessionNotification');
EXCEPTION
 WHEN duplicate_object THEN null;
END $$;
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "email"."email" (
	"id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"sub" uuid NOT NULL,
	"sent_at" timestamp NOT NULL,
	"template" "templates_enum" NOT NULL
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "email"."preference" (
	"sub" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"email_address" varchar NOT NULL,
	"allow_account_reset" boolean DEFAULT true NOT NULL,
	"allow_session_notification" boolean DEFAULT false NOT NULL
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "email"."resend_email" (
	"id" uuid PRIMARY KEY NOT NULL,
	"resend_id" varchar NOT NULL,
	CONSTRAINT "resend_email_resend_id_unique" UNIQUE("resend_id")
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "email"."verification" (
	"id" uuid PRIMARY KEY NOT NULL,
	"verified_at" timestamp NOT NULL
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "email"."verification_request" (
	"id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"sub" uuid NOT NULL,
	"expires_at" timestamp NOT NULL,
	"code_verifier" char(48) NOT NULL
);
--> statement-breakpoint
DO $$ BEGIN
 ALTER TABLE "email"."email" ADD CONSTRAINT "email_sub_preference_sub_fk" FOREIGN KEY ("sub") REFERENCES "email"."preference"("sub") ON DELETE no action ON UPDATE no action;
EXCEPTION
 WHEN duplicate_object THEN null;
END $$;
--> statement-breakpoint
DO $$ BEGIN
 ALTER TABLE "email"."resend_email" ADD CONSTRAINT "resend_email_id_email_id_fk" FOREIGN KEY ("id") REFERENCES "email"."email"("id") ON DELETE cascade ON UPDATE no action;
EXCEPTION
 WHEN duplicate_object THEN null;
END $$;
--> statement-breakpoint
DO $$ BEGIN
 ALTER TABLE "email"."verification" ADD CONSTRAINT "verification_id_verification_request_id_fk" FOREIGN KEY ("id") REFERENCES "email"."verification_request"("id") ON DELETE cascade ON UPDATE no action;
EXCEPTION
 WHEN duplicate_object THEN null;
END $$;
--> statement-breakpoint
DO $$ BEGIN
 ALTER TABLE "email"."verification_request" ADD CONSTRAINT "verification_request_sub_preference_sub_fk" FOREIGN KEY ("sub") REFERENCES "email"."preference"("sub") ON DELETE no action ON UPDATE no action;
EXCEPTION
 WHEN duplicate_object THEN null;
END $$;
