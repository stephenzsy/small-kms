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
 * @interface ExchangePKCS12Request
 */
export interface ExchangePKCS12Request {
    /**
     * JWE encrypted private key in JWK
     * @type {string}
     * @memberof ExchangePKCS12Request
     */
    payload: string;
    /**
     * Use legacy PKCS12 cipher
     * @type {boolean}
     * @memberof ExchangePKCS12Request
     */
    legacy?: boolean;
    /**
     * Encrypt the PKCS12 file with a generated password
     * @type {boolean}
     * @memberof ExchangePKCS12Request
     */
    passwordProtected: boolean;
}

/**
 * Check if a given object implements the ExchangePKCS12Request interface.
 */
export function instanceOfExchangePKCS12Request(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "payload" in value;
    isInstance = isInstance && "passwordProtected" in value;

    return isInstance;
}

export function ExchangePKCS12RequestFromJSON(json: any): ExchangePKCS12Request {
    return ExchangePKCS12RequestFromJSONTyped(json, false);
}

export function ExchangePKCS12RequestFromJSONTyped(json: any, ignoreDiscriminator: boolean): ExchangePKCS12Request {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'payload': json['payload'],
        'legacy': !exists(json, 'legacy') ? undefined : json['legacy'],
        'passwordProtected': json['passwordProtected'],
    };
}

export function ExchangePKCS12RequestToJSON(value?: ExchangePKCS12Request | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'payload': value.payload,
        'legacy': value.legacy,
        'passwordProtected': value.passwordProtected,
    };
}

