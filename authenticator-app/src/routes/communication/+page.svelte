<script lang="ts">
  import CardWithIcon from "$lib/CardWithIcon.svelte";
  import Email from "$lib/login/icons/Email.svelte";
  import Input from "$lib/utils/Input.svelte";
  import type { PageServerData } from "./$types";

  export let data: PageServerData;
</script>

<main class="sm:w-[420px]">
  <CardWithIcon>
    <svelte:fragment slot="icon">
      <Email />
    </svelte:fragment>
    <h1 class="relative z-100">Communication</h1>
    <h4>Update communication data and preferences</h4>
    <section class="mt-2">
      <div>
        <h2>
          Email Address
          {#if data.verified}
          <span class="pill pill-sm pill-green">
            Verified
          </span>
          {:else}
          <span class="pill pill-red">
            Verification required
          </span>
          {/if}
        </h2>
        <div class="flex flex-col space-y-1">
          <Input
            class="grow"
            inputClass="rounded-br-sm"
            id="email"
            label="Email"
            value={data.email}
            autocomplete="email"
            placeholder="abc@example.com"
            onInput={() => {}}
          />
          <button class="btn rounded-b-lg rounded-t-sm border self-end">Resend verification</button>
        </div>
      </div>
    </section>
    <section class="mt-2">
      <h2>Preferences</h2>
      <ul>
        <li>
          <input id="account-reset" type="checkbox" class="rounded-sm text-yellow-500" checked />
          <label for="account-reset">Use email to reset account</label>
        </li>
        <li>
          <input id="notify-session" type="checkbox" class="rounded-sm text-yellow-500" checked />
          <label for="notify-session">Notify about new sessions</label>
        </li>
      </ul>
    </section>
  </CardWithIcon>
  <p class="card mt-2 card-blue">
    The information you enter here, will <em>never be shared with third parties</em>. The email entered here is only used for communication relating to you account, such as account recovery.
  </p>
</main>
