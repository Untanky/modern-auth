<script lang="ts">
  import Input from "../utils/Input.svelte";

  export let submit: () => void;
</script>

<form class="flex flex-col flex-1 space-y-1">
  <h2 class="text-xl">
    Profile
  </h2>
  <p>
    When you login to 3rd party apps through this account, you can select to share this information with the app. You may complete your profile:
  </p>
  <div class="grid grid-cols-2 py-4 gap-2">
    <h3 class="text-lg col-span-2">Personal information:</h3>
    <Input
      type="text"
      value=""
      placeholder="Nomen"
      autocomplete="given-name"
      onInput={() => {}}
      id="firstname"
      label="firstname"
    />
    <Input
      type="text"
      value=""
      placeholder="Nescio"
      autocomplete="family-name"
      onInput={() => {}}
      id="lastname"
      label="lastname"
    />
    <Input
      class="col-span-2"
      type="text"
      value=""
      placeholder="Your Username"
      autocomplete="username"
      onInput={() => {}}
      id="username"
      label="username"
    />
    <Input
      class="col-span-2"
      type="text"
      value=""
      placeholder="nomen.nescio@email.com"
      autocomplete="email"
      onInput={() => {}}
      id="email"
      label="email"
    />
    <h3 class="text-lg">Address:</h3>
    <Input
      class="col-span-2"
      type="text"
      value=""
      placeholder="1234 Example Street"
      autocomplete="street-address"
      onInput={() => {}}
      id="street"
      label="street"
    />
    <Input
      type="text"
      value=""
      placeholder="12345"
      autocomplete="postal-code"
      onInput={() => {}}
      id="zip"
      label="zip"
    />
    <Input
      type="text"
      value=""
      placeholder="Example City"
      autocomplete="address-level2"
      onInput={() => {}}
      id="city"
      label="city"
    />
    <Input
      class="col-span-2"
      type="text"
      value=""
      placeholder="Example Country"
      autocomplete="country-name"
      onInput={() => {}}
      id="country"
      label="country"
    />
  </div>
  <div  class="self-end mt-2 space-x-2">
    <button type="button" on:click={submit} class="btn btn-link">
      Skip
    </button>
    <button type="button" on:click={submit} class="btn btn-primary">
      Register
    </button>
  </div>
</form>
