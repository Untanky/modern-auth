variable VERSION {
    default = "0.1-beta"
}

variable IS_LATEST {
    default = true
}

variable REF {
    default = ""
}

group default {
    targets = [
        "authenticator-app",
        "oauth2",
        "webauthn",
        "passwords",
    ]
}

target authenticator-app {
    dockerfile = "Dockerfile"
    context = "authenticator-app"
    labels = {
        "org.opencontainers.image.name" = "modern-auth/authenticator-app"
        "org.opencontainers.image.description" = "App containing the frontend, email and profile service of ModernAuth"
        "org.opencontainers.image.url" = "https://github.com/Untanky/modern-auth"
        "org.opencontainers.image.source" = "https://github.com/Untanky/modern-auth"
        "org.opencontainers.image.revision" = "${REF}"
        "org.opencontainers.image.version=" = "${VERSION}"
    }
    tags = [
        // "untanky/authenticator-app:${VERSION}",
        IS_LATEST ? "ghcr.io/untanky/authenticator-app:latest" : "",
        "ghcr.io/untanky/authenticator-app:${VERSION}",
    ]
    cache-from = [
        "type=gha"
    ]
    cache-to = [
        "type=gha,mode=max"
    ]
    platforms = [
        "linux/amd64"
    ]
}

target oauth2 {
    dockerfile = "Dockerfile"
    context = "apps/oauth2"
    labels = {
        "org.opencontainers.image.name" = "modern-auth/oauth2-service"
        "org.opencontainers.image.description" = "Service for handling OAuth2.1"
        "org.opencontainers.image.url" = "https://github.com/Untanky/modern-auth"
        "org.opencontainers.image.source" = "https://github.com/Untanky/modern-auth"
        "org.opencontainers.image.revision" = "${REF}"
        "org.opencontainers.image.version=" = "${VERSION}"
    }
    tags = [
        // "untanky/authenticator-app:${VERSION}",
        IS_LATEST ? "ghcr.io/untanky/oauth2-service:latest" : "",
        "ghcr.io/untanky/oauth2-service:${VERSION}",
    ]
    cache-from = [
        "type=gha"
    ]
    cache-to = [
        "type=gha,mode=max"
    ]
    platforms = [
        "linux/amd64"
    ]
}

target webauthn {
    dockerfile = "Dockerfile"
    context = "apps/webauthn"
    labels = {
        "org.opencontainers.image.name" = "modern-auth/webauthn-service"
        "org.opencontainers.image.description" = "Service for handling WebAuthn"
        "org.opencontainers.image.url" = "https://github.com/Untanky/modern-auth"
        "org.opencontainers.image.source" = "https://github.com/Untanky/modern-auth"
        "org.opencontainers.image.revision" = "${REF}"
        "org.opencontainers.image.version=" = "${VERSION}"
    }
    tags = [
        // "untanky/authenticator-app:${VERSION}",
        IS_LATEST ? "ghcr.io/untanky/webauthn-service:latest" : "",
        "ghcr.io/untanky/webauthn-service:${VERSION}",
    ]
    cache-from = [
        "type=gha"
    ]
    cache-to = [
        "type=gha,mode=max"
    ]
    platforms = [
        "linux/amd64"
    ]
}

target passwords {
    dockerfile = "Dockerfile"
    context = "apps/passwords"
    labels = {
        "org.opencontainers.image.name" = "modern-auth/passwords-service"
        "org.opencontainers.image.description" = "Service for handling passwords over WebAuthn"
        "org.opencontainers.image.url" = "https://github.com/Untanky/modern-auth"
        "org.opencontainers.image.source" = "https://github.com/Untanky/modern-auth"
        "org.opencontainers.image.revision" = "${REF}"
        "org.opencontainers.image.version=" = "${VERSION}"
    }
    tags = [
        // "untanky/authenticator-app:${VERSION}",
        IS_LATEST ? "ghcr.io/untanky/passwords-service:latest" : "",
        "ghcr.io/untanky/passwords-service:${VERSION}",
    ]
    cache-from = [
        "type=gha"
    ]
    cache-to = [
        "type=gha,mode=max"
    ]
    platforms = [
        "linux/amd64"
    ]
}
