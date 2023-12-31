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
import type { CertificatePendingAcmeChallenge } from './CertificatePendingAcmeChallenge';
import {
    CertificatePendingAcmeChallengeFromJSON,
    CertificatePendingAcmeChallengeFromJSONTyped,
    CertificatePendingAcmeChallengeToJSON,
} from './CertificatePendingAcmeChallenge';

/**
 * 
 * @export
 * @interface CertificatePendingAcmeAuthorization
 */
export interface CertificatePendingAcmeAuthorization {
    /**
     * URL to the authorization
     * @type {string}
     * @memberof CertificatePendingAcmeAuthorization
     */
    url: string;
    /**
     * 
     * @type {string}
     * @memberof CertificatePendingAcmeAuthorization
     */
    status: string;
    /**
     * 
     * @type {Array<CertificatePendingAcmeChallenge>}
     * @memberof CertificatePendingAcmeAuthorization
     */
    challenges?: Array<CertificatePendingAcmeChallenge>;
}

/**
 * Check if a given object implements the CertificatePendingAcmeAuthorization interface.
 */
export function instanceOfCertificatePendingAcmeAuthorization(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "url" in value;
    isInstance = isInstance && "status" in value;

    return isInstance;
}

export function CertificatePendingAcmeAuthorizationFromJSON(json: any): CertificatePendingAcmeAuthorization {
    return CertificatePendingAcmeAuthorizationFromJSONTyped(json, false);
}

export function CertificatePendingAcmeAuthorizationFromJSONTyped(json: any, ignoreDiscriminator: boolean): CertificatePendingAcmeAuthorization {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'url': json['url'],
        'status': json['status'],
        'challenges': !exists(json, 'challenges') ? undefined : ((json['challenges'] as Array<any>).map(CertificatePendingAcmeChallengeFromJSON)),
    };
}

export function CertificatePendingAcmeAuthorizationToJSON(value?: CertificatePendingAcmeAuthorization | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'url': value.url,
        'status': value.status,
        'challenges': value.challenges === undefined ? undefined : ((value.challenges as Array<any>).map(CertificatePendingAcmeChallengeToJSON)),
    };
}

