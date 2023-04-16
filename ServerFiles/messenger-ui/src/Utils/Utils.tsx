import CryptoJS from "crypto-js";

const base64_iv  = 'fb4c5e213749eddadf1e22d723eaf207';
const iv  = CryptoJS.enc.Hex.parse(base64_iv);

var Key:CryptoJS.lib.WordArray = null!

function setKey(key:CryptoJS.lib.WordArray){
    Key = key
}

function decryptText<T>(encryptedText:string,key:CryptoJS.lib.WordArray):T{
    key = key
    let plain = CryptoJS.AES.decrypt(encryptedText,key,
                                        { iv: iv,
                                            mode: CryptoJS.mode.CBC,
                                            padding: CryptoJS.pad.Pkcs7})
    let result:T = JSON.parse(plain.toString(CryptoJS.enc.Utf8))
    return result
}
function encryptObject(object:any,key:CryptoJS.lib.WordArray){
    key = key
    if(typeof object != "string" )
       object =  JSON.stringify(object)

    let encryptedObject = CryptoJS.AES.encrypt(object,
                                        key,
                                        { iv: iv,
                                          mode: CryptoJS.mode.CBC,
                                          padding: CryptoJS.pad.Pkcs7})      
    return encryptedObject.toString()
}

export {encryptObject,decryptText,setKey,Key};