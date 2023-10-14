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
 * @interface ServiceConfigAcrInfo
 */
export interface ServiceConfigAcrInfo {
    /**
     * 
     * @type {string}
     * @memberof ServiceConfigAcrInfo
     */
    loginServer: string;
    /**
     * 
     * @type {string}
     * @memberof ServiceConfigAcrInfo
     */
    name: string;
    /**
     * 
     * @type {string}
     * @memberof ServiceConfigAcrInfo
     */
    armResourceId: string;
}

/**
 * Check if a given object implements the ServiceConfigAcrInfo interface.
 */
export function instanceOfServiceConfigAcrInfo(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "loginServer" in value;
    isInstance = isInstance && "name" in value;
    isInstance = isInstance && "armResourceId" in value;

    return isInstance;
}

export function ServiceConfigAcrInfoFromJSON(json: any): ServiceConfigAcrInfo {
    return ServiceConfigAcrInfoFromJSONTyped(json, false);
}

export function ServiceConfigAcrInfoFromJSONTyped(json: any, ignoreDiscriminator: boolean): ServiceConfigAcrInfo {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'loginServer': json['loginServer'],
        'name': json['name'],
        'armResourceId': json['armResourceId'],
    };
}

export function ServiceConfigAcrInfoToJSON(value?: ServiceConfigAcrInfo | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'loginServer': value.loginServer,
        'name': value.name,
        'armResourceId': value.armResourceId,
    };
}

