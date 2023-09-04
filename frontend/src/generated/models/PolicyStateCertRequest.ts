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
/**
 * 
 * @export
 * @interface PolicyStateCertRequest
 */
export interface PolicyStateCertRequest {
    /**
     * 
     * @type {string}
     * @memberof PolicyStateCertRequest
     */
    lastCertId: string;
    /**
     * 
     * @type {Date}
     * @memberof PolicyStateCertRequest
     */
    lastCertIssued: Date;
    /**
     * 
     * @type {Date}
     * @memberof PolicyStateCertRequest
     */
    lastCertExpires: Date;
    /**
     * 
     * @type {string}
     * @memberof PolicyStateCertRequest
     */
    lastAction: string;
}

/**
 * Check if a given object implements the PolicyStateCertRequest interface.
 */
export function instanceOfPolicyStateCertRequest(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "lastCertId" in value;
    isInstance = isInstance && "lastCertIssued" in value;
    isInstance = isInstance && "lastCertExpires" in value;
    isInstance = isInstance && "lastAction" in value;

    return isInstance;
}

export function PolicyStateCertRequestFromJSON(json: any): PolicyStateCertRequest {
    return PolicyStateCertRequestFromJSONTyped(json, false);
}

export function PolicyStateCertRequestFromJSONTyped(json: any, ignoreDiscriminator: boolean): PolicyStateCertRequest {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'lastCertId': json['lastCertId'],
        'lastCertIssued': (new Date(json['lastCertIssued'])),
        'lastCertExpires': (new Date(json['lastCertExpires'])),
        'lastAction': json['lastAction'],
    };
}

export function PolicyStateCertRequestToJSON(value?: PolicyStateCertRequest | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'lastCertId': value.lastCertId,
        'lastCertIssued': (value.lastCertIssued.toISOString()),
        'lastCertExpires': (value.lastCertExpires.toISOString()),
        'lastAction': value.lastAction,
    };
}

