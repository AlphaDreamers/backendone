"use client";

import {ethers} from "ethers";
import {
    Card,
    CardContent,
    CardDescription,
    CardFooter,
    CardHeader,
    CardTitle,
} from "@/components/ui/card"



export  default  function Wallet() {

    const wallet = ethers.Wallet.createRandom();
    const mnemonic =  wallet.mnemonic?.phrase
    return <>
        <Card>
            <CardHeader>
                <CardTitle>Card Title</CardTitle>
                <CardDescription>Card Description</CardDescription>
            </CardHeader>
            <CardContent>
                <p>{mnemonic}</p>
            </CardContent>
            <CardFooter>
                <p>Card Footer</p>
            </CardFooter>
        </Card>

    </>
}