# confidential-secrets-demo

A minimal Tinfoil Containers workload that receives a secret from a
customer-run keyserver — released only to this repo's attested releases —
and proves possession without ever disclosing it.

At boot, the enclave fetches `DEMO_SECRET` from the keyserver pinned in
`tinfoil-config.yml` (`vault-url`) and injects it as an env var. The workload
exposes:

- `GET /` — `{"secret_present": true}`
- `GET /prove?challenge=<text>` — hex `HMAC-SHA256(DEMO_SECRET, text)`; anyone
  who knows the secret can recompute it, nobody else learns anything.

## Deploy

1. Run a [keyserver](https://github.com/tinfoilsh/keyserver) at a domain you
   control, and set that domain as `vault-url` in `tinfoil-config.yml`.
2. Trigger the **Tinfoil Release** workflow with a version — it builds the
   image, pins its digest into the config, tags, measures, and publishes a
   Sigstore-attested release.
3. Pin the release in the keyserver's policy and store the secret:

   ```yaml
   workloads:
     demo:
       repo: tinfoilsh/confidential-secrets-demo
       tag: <the released version>
       domain: <your deployment's domain>
       secrets:
         DEMO_SECRET: {path: workloads/demo/secret, field: value}
   ```

4. Deploy this repo@tag on Tinfoil Containers, then verify:

   ```bash
   curl https://<deployment>/prove?challenge=hello
   ```

The release is fail-closed end to end: the keyserver releases only to a fresh,
hardware-attested enclave running this exact release, and the container never
starts with the secret missing.
