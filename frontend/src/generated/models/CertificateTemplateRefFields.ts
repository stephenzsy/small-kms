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
 * @interface CertificateTemplateRefFields
 */
export interface CertificateTemplateRefFields {
    /**
     * Common name
     * @type {string}
     * @memberof CertificateTemplateRefFields
     */
    subjectCommonName: string;
    /**
     * 
     * @type {string}
     * @memberof CertificateTemplateRefFields
     */
    linkTo?: string;
}

/**
 * Check if a given object implements the CertificateTemplateRefFields interface.
 */
export function instanceOfCertificateTemplateRefFields(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "subjectCommonName" in value;

    return isInstance;
}

export function CertificateTemplateRefFieldsFromJSON(json: any): CertificateTemplateRefFields {
    return CertificateTemplateRefFieldsFromJSONTyped(json, false);
}

export function CertificateTemplateRefFieldsFromJSONTyped(json: any, ignoreDiscriminator: boolean): CertificateTemplateRefFields {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'subjectCommonName': json['subjectCommonName'],
        'linkTo': !exists(json, 'linkTo') ? undefined : json['linkTo'],
    };
}

export function CertificateTemplateRefFieldsToJSON(value?: CertificateTemplateRefFields | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'subjectCommonName': value.subjectCommonName,
        'linkTo': value.linkTo,
    };
}
