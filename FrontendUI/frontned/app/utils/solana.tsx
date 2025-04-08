const { Keypair } = require("@solana/web3.js");
const bip39 = require("bip39");
const { derivePath } = require("ed25519-hd-key");
const bs58 = require("bs58");

async function generateWallet() {
    const mnemonic = bip39.generateMnemonic(128); // 128 bits for 12 words, 256 for 24 words

    const seed = await bip39.mnemonicToSeed(mnemonic);

    const derivationPath = "m/44'/501'/0'/0'"; 
    const derivedSeed = derivePath(derivationPath, seed.toString("hex")).key;
    const wallet = Keypair.fromSecretKey(derivedSeed);

    console.log("📝 Mnemonic:", mnemonic);
    console.log("✅ Public Key:", wallet.publicKey.toBase58());
    console.log("🔑 Private Key:", bs58.encode(wallet.secretKey));

}