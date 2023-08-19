<script lang="ts">
  import Input from "$lib/utils/Input.svelte";

  export let disabled = false;
</script>

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
    disabled={disabled}
  />
  <Input
    type="text"
    value=""
    placeholder="Nescio"
    autocomplete="family-name"
    onInput={() => {}}
    id="lastname"
    label="lastname"
    disabled={disabled}
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
    disabled={disabled}
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
    disabled={disabled}
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
    disabled={disabled}
  />
  <Input
    type="text"
    value=""
    placeholder="12345"
    autocomplete="postal-code"
    onInput={() => {}}
    id="zip"
    label="zip"
    disabled={disabled}
  />
  <Input
    type="text"
    value=""
    placeholder="Example City"
    autocomplete="address-level2"
    onInput={() => {}}
    id="city"
    label="city"
    disabled={disabled}
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
    disabled={disabled}
  />
</div>