<script lang="ts">
    import { Button, Container, Html, Link, Text } from 'svelte-email';
    export let name: string;
    export let verificationLink: string = 'http://localhost:5173/v1/verification';
</script>

<Html lang="en">
    <Container>
        <Text>
            Hello {name},
        </Text>
        <Text>
            your email has been used to register for ModernAuth. To verify you email please click on the button.
        </Text>
        <Button style={{ backgroundColor: 'rgb(250, 204, 21)', padding: '8px 16px', borderRadius: '16px' }} href={verificationLink}>
            Verify
        </Button>
        <Text>
            If you can click the button, copy the following URL into your browser:
            <Link href={verificationLink}>
                {verificationLink}
            </Link>
        </Text>
    </Container>
</Html>