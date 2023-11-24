/* tslint:disable */
/* eslint-disable */
/**
 * Cryptocat API
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: 0.1.3
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */

import { exists, mapValues } from '../runtime';
/**
 * 
 * @export
 * @interface ExchangePKCS12Result
 */
export interface ExchangePKCS12Result {
    /**
     * JWE encrypted PKCS12 file, encrypted with the symmetric key from the request
     * @type {string}
     * @memberof ExchangePKCS12Result
     */
    payload: string;
    /**
     * Password used to encrypt the PKCS12 file
     * @type {string}
     * @memberof ExchangePKCS12Result
     */
    password: string;
}

/**
 * Check if a given object implements the ExchangePKCS12Result interface.
 */
export function instanceOfExchangePKCS12Result(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "payload" in value;
    isInstance = isInstance && "password" in value;

    return isInstance;
}

export function ExchangePKCS12ResultFromJSON(json: any): ExchangePKCS12Result {
    return ExchangePKCS12ResultFromJSONTyped(json, false);
}

export function ExchangePKCS12ResultFromJSONTyped(json: any, ignoreDiscriminator: boolean): ExchangePKCS12Result {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'payload': json['payload'],
        'password': json['password'],
    };
}

export function ExchangePKCS12ResultToJSON(value?: ExchangePKCS12Result | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'payload': value.payload,
        'password': value.password,
    };
}

