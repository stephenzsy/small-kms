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
 * @interface CertificatePendingAcmeChallenge
 */
export interface CertificatePendingAcmeChallenge {
    /**
     * 
     * @type {string}
     * @memberof CertificatePendingAcmeChallenge
     */
    type: string;
    /**
     * URL to the challenge
     * @type {string}
     * @memberof CertificatePendingAcmeChallenge
     */
    url: string;
    /**
     * 
     * @type {string}
     * @memberof CertificatePendingAcmeChallenge
     */
    dnsRecord?: string;
}

/**
 * Check if a given object implements the CertificatePendingAcmeChallenge interface.
 */
export function instanceOfCertificatePendingAcmeChallenge(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "type" in value;
    isInstance = isInstance && "url" in value;

    return isInstance;
}

export function CertificatePendingAcmeChallengeFromJSON(json: any): CertificatePendingAcmeChallenge {
    return CertificatePendingAcmeChallengeFromJSONTyped(json, false);
}

export function CertificatePendingAcmeChallengeFromJSONTyped(json: any, ignoreDiscriminator: boolean): CertificatePendingAcmeChallenge {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'type': json['type'],
        'url': json['url'],
        'dnsRecord': !exists(json, 'dnsRecord') ? undefined : json['dnsRecord'],
    };
}

export function CertificatePendingAcmeChallengeToJSON(value?: CertificatePendingAcmeChallenge | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'type': value.type,
        'url': value.url,
        'dnsRecord': value.dnsRecord,
    };
}

