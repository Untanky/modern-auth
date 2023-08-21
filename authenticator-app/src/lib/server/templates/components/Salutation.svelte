<script lang="ts">
  import { Text } from "svelte-email";

    export let userName: string = '';
    export let firstName: string = '';
    export let lastName: string = '';

    $: fullName = `${firstName} ${lastName}`.trim();
    $: name = fullName.length > 0 ? `Mrs/Mr. ${fullName}` : userName;
</script>

<Text>
    Hello {name},
</Text>