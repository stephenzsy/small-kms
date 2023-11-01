/* tslint:disable */
/* eslint-disable */
/**
 * Small KMS Admin API
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: 0.1.1
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
 * @interface AuthResult
 */
export interface AuthResult {
    /**
     * 
     * @type {string}
     * @memberof AuthResult
     */
    accessToken: string;
}

/**
 * Check if a given object implements the AuthResult interface.
 */
export function instanceOfAuthResult(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "accessToken" in value;

    return isInstance;
}

export function AuthResultFromJSON(json: any): AuthResult {
    return AuthResultFromJSONTyped(json, false);
}

export function AuthResultFromJSONTyped(json: any, ignoreDiscriminator: boolean): AuthResult {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'accessToken': json['accessToken'],
    };
}

export function AuthResultToJSON(value?: AuthResult | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'accessToken': value.accessToken,
    };
}
