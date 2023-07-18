<script lang="ts">
  import { initiateAuthentication } from "./Authentication";
  import CreateCredential from "./CreateCredential.svelte";
  import GetCredential from "./GetCredential.svelte";
  import UserIdForm from "./UserIdForm.svelte";

  type State = 'userId' | 'createCredential' | 'getCredential' | 'success';
  
  let state: State = 'userId';
  let loading = false;
  let error: Error | null = null;
  let userId: string | null = null;

  let credOps: CredentialCreationOptions | null = null;

  const onInitiateAuthentication = (): void => {
    const inputElement = document.getElementById('user-id') as HTMLInputElement;
    const inputValue = inputElement.value;

    userId = inputValue;
    loading = true;

    initiateAuthentication(inputValue)
      .then((ops) => {
        credOps = ops;
        state = 'createCredential';
      })
      .catch((err) => {
        error = err;
      })
      .finally(() => {
        loading = false;
      });
  };

  const onCreateCredential = (): void => {
    navigator.credentials.create({
      publicKey: {
        ...credOps.publicKey,
        challenge: parseStringToUint8Array(credOps.publicKey.challenge as unknown as string),
        rp: {
          id: credOps.publicKey.rp.id,
          name: credOps.publicKey.rp.name,
        },
        user: {
          id: parseStringToUint8Array(credOps.publicKey.user.id as unknown as string),
          name: credOps.publicKey.user.name,
          displayName: credOps.publicKey.user.displayName,
        },
      },
    });
  };

  const onGetCredential = (): void => {
    console.log('create credential');
  };

  // parse string to uint8 array
  const parseStringToUint8Array = (str: string): Uint8Array => {
    const arr = new Uint8Array(str.length);
    for (let i = 0; i < str.length; i++) {
      arr[i] = str.charCodeAt(i);
    }
    return arr;
  };

</script>

<div class="flex w-screen h-screen justify-center items-center">
  <main class="w-[420px] card">
    <h1 class="text-stone-950 dark:text-stone-50 text-xl font-medium">
      Login
    </h1>
    {#if state === 'userId'}
    <UserIdForm submit={onInitiateAuthentication} />
    {:else if state === 'createCredential'}
    <CreateCredential userId={userId} submit={onCreateCredential} />
    {:else if state === 'getCredential'}
    <GetCredential userId={userId} submit={onGetCredential} />
    {/if}
  </main>
</div>
