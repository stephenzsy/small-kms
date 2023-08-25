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
import type { CertificateCategory } from './CertificateCategory';
import {
    CertificateCategoryFromJSON,
    CertificateCategoryFromJSONTyped,
    CertificateCategoryToJSON,
} from './CertificateCategory';
import type { CertificateSubject } from './CertificateSubject';
import {
    CertificateSubjectFromJSON,
    CertificateSubjectFromJSONTyped,
    CertificateSubjectToJSON,
} from './CertificateSubject';
import type { CreateCertificateOptions } from './CreateCertificateOptions';
import {
    CreateCertificateOptionsFromJSON,
    CreateCertificateOptionsFromJSONTyped,
    CreateCertificateOptionsToJSON,
} from './CreateCertificateOptions';

/**
 * 
 * @export
 * @interface CreateCertificateParameters
 */
export interface CreateCertificateParameters {
    /**
     * 
     * @type {CertificateCategory}
     * @memberof CreateCertificateParameters
     */
    category: CertificateCategory;
    /**
     * 
     * @type {string}
     * @memberof CreateCertificateParameters
     */
    name: string;
    /**
     * 
     * @type {CertificateSubject}
     * @memberof CreateCertificateParameters
     */
    subject: CertificateSubject;
    /**
     * 
     * @type {string}
     * @memberof CreateCertificateParameters
     */
    validity?: string;
    /**
     * 
     * @type {string}
     * @memberof CreateCertificateParameters
     */
    kty?: CreateCertificateParametersKtyEnum;
    /**
     * 
     * @type {number}
     * @memberof CreateCertificateParameters
     */
    size?: CreateCertificateParametersSizeEnum;
    /**
     * 
     * @type {string}
     * @memberof CreateCertificateParameters
     */
    curve?: CreateCertificateParametersCurveEnum;
    /**
     * 
     * @type {string}
     * @memberof CreateCertificateParameters
     */
    issuer?: string;
    /**
     * 
     * @type {CreateCertificateOptions}
     * @memberof CreateCertificateParameters
     */
    options?: CreateCertificateOptions;
}


/**
 * @export
 */
export const CreateCertificateParametersKtyEnum = {
    Rsa: 'RSA',
    Ec: 'EC'
} as const;
export type CreateCertificateParametersKtyEnum = typeof CreateCertificateParametersKtyEnum[keyof typeof CreateCertificateParametersKtyEnum];

/**
 * @export
 */
export const CreateCertificateParametersSizeEnum = {
    NUMBER_2048: 2048,
    NUMBER_3072: 3072,
    NUMBER_4096: 4096
} as const;
export type CreateCertificateParametersSizeEnum = typeof CreateCertificateParametersSizeEnum[keyof typeof CreateCertificateParametersSizeEnum];

/**
 * @export
 */
export const CreateCertificateParametersCurveEnum = {
    _256: 'P-256',
    _384: 'P-384'
} as const;
export type CreateCertificateParametersCurveEnum = typeof CreateCertificateParametersCurveEnum[keyof typeof CreateCertificateParametersCurveEnum];


/**
 * Check if a given object implements the CreateCertificateParameters interface.
 */
export function instanceOfCreateCertificateParameters(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "category" in value;
    isInstance = isInstance && "name" in value;
    isInstance = isInstance && "subject" in value;

    return isInstance;
}

export function CreateCertificateParametersFromJSON(json: any): CreateCertificateParameters {
    return CreateCertificateParametersFromJSONTyped(json, false);
}

export function CreateCertificateParametersFromJSONTyped(json: any, ignoreDiscriminator: boolean): CreateCertificateParameters {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'category': CertificateCategoryFromJSON(json['category']),
        'name': json['name'],
        'subject': CertificateSubjectFromJSON(json['subject']),
        'validity': !exists(json, 'validity') ? undefined : json['validity'],
        'kty': !exists(json, 'kty') ? undefined : json['kty'],
        'size': !exists(json, 'size') ? undefined : json['size'],
        'curve': !exists(json, 'curve') ? undefined : json['curve'],
        'issuer': !exists(json, 'issuer') ? undefined : json['issuer'],
        'options': !exists(json, 'options') ? undefined : CreateCertificateOptionsFromJSON(json['options']),
    };
}

export function CreateCertificateParametersToJSON(value?: CreateCertificateParameters | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'category': CertificateCategoryToJSON(value.category),
        'name': value.name,
        'subject': CertificateSubjectToJSON(value.subject),
        'validity': value.validity,
        'kty': value.kty,
        'size': value.size,
        'curve': value.curve,
        'issuer': value.issuer,
        'options': CreateCertificateOptionsToJSON(value.options),
    };
}

