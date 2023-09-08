/* tslint:disable */
/* eslint-disable */
/**
 * Small KMS Admin API
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: 0.1.0
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */

import { exists, mapValues } from '../runtime';
import type { PolicyType } from './PolicyType';
import {
    PolicyTypeFromJSON,
    PolicyTypeFromJSONTyped,
    PolicyTypeToJSON,
} from './PolicyType';

/**
 * 
 * @export
 * @interface PolicyRefParameters
 */
export interface PolicyRefParameters {
    /**
     * 
     * @type {PolicyType}
     * @memberof PolicyRefParameters
     */
    policyType: PolicyType;
}

/**
 * Check if a given object implements the PolicyRefParameters interface.
 */
export function instanceOfPolicyRefParameters(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "policyType" in value;

    return isInstance;
}

export function PolicyRefParametersFromJSON(json: any): PolicyRefParameters {
    return PolicyRefParametersFromJSONTyped(json, false);
}

export function PolicyRefParametersFromJSONTyped(json: any, ignoreDiscriminator: boolean): PolicyRefParameters {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'policyType': PolicyTypeFromJSON(json['policyType']),
    };
}

export function PolicyRefParametersToJSON(value?: PolicyRefParameters | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'policyType': PolicyTypeToJSON(value.policyType),
    };
}
