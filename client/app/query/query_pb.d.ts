// package: query
// file: query/query.proto

import * as jspb from "google-protobuf";

export class EnvelopeToSign extends jspb.Message {
  getEnvelope(): Uint8Array | string;
  getEnvelope_asU8(): Uint8Array;
  getEnvelope_asB64(): string;
  setEnvelope(value: Uint8Array | string): void;

  getPollid(): number;
  setPollid(value: number): void;

  hasToken(): boolean;
  clearToken(): void;
  getToken(): VoteToken | undefined;
  setToken(value?: VoteToken): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnvelopeToSign.AsObject;
  static toObject(includeInstance: boolean, msg: EnvelopeToSign): EnvelopeToSign.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: EnvelopeToSign, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EnvelopeToSign;
  static deserializeBinaryFromReader(message: EnvelopeToSign, reader: jspb.BinaryReader): EnvelopeToSign;
}

export namespace EnvelopeToSign {
  export type AsObject = {
    envelope: Uint8Array | string,
    pollid: number,
    token?: VoteToken.AsObject,
  }
}

export class PollWithPublicKey extends jspb.Message {
  hasKey(): boolean;
  clearKey(): void;
  getKey(): PublicKey | undefined;
  setKey(value?: PublicKey): void;

  hasPoll(): boolean;
  clearPoll(): void;
  getPoll(): PollSchema | undefined;
  setPoll(value?: PollSchema): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PollWithPublicKey.AsObject;
  static toObject(includeInstance: boolean, msg: PollWithPublicKey): PollWithPublicKey.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: PollWithPublicKey, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PollWithPublicKey;
  static deserializeBinaryFromReader(message: PollWithPublicKey, reader: jspb.BinaryReader): PollWithPublicKey;
}

export namespace PollWithPublicKey {
  export type AsObject = {
    key?: PublicKey.AsObject,
    poll?: PollSchema.AsObject,
  }
}

export class GetPollRequest extends jspb.Message {
  getPollid(): number;
  setPollid(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetPollRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetPollRequest): GetPollRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetPollRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetPollRequest;
  static deserializeBinaryFromReader(message: GetPollRequest, reader: jspb.BinaryReader): GetPollRequest;
}

export namespace GetPollRequest {
  export type AsObject = {
    pollid: number,
  }
}

export class PollAnswer extends jspb.Message {
  hasAnswers(): boolean;
  clearAnswers(): void;
  getAnswers(): PollSchema | undefined;
  setAnswers(value?: PollSchema): void;

  hasSign(): boolean;
  clearSign(): void;
  getSign(): RSASignature | undefined;
  setSign(value?: RSASignature): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PollAnswer.AsObject;
  static toObject(includeInstance: boolean, msg: PollAnswer): PollAnswer.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: PollAnswer, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PollAnswer;
  static deserializeBinaryFromReader(message: PollAnswer, reader: jspb.BinaryReader): PollAnswer;
}

export namespace PollAnswer {
  export type AsObject = {
    answers?: PollSchema.AsObject,
    sign?: RSASignature.AsObject,
  }
}

export class PollSchema extends jspb.Message {
  clearQuestionsList(): void;
  getQuestionsList(): Array<PollSchema.QA>;
  setQuestionsList(value: Array<PollSchema.QA>): void;
  addQuestions(value?: PollSchema.QA, index?: number): PollSchema.QA;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PollSchema.AsObject;
  static toObject(includeInstance: boolean, msg: PollSchema): PollSchema.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: PollSchema, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PollSchema;
  static deserializeBinaryFromReader(message: PollSchema, reader: jspb.BinaryReader): PollSchema;
}

export namespace PollSchema {
  export type AsObject = {
    questionsList: Array<PollSchema.QA.AsObject>,
  }

  export class QA extends jspb.Message {
    getQuestion(): string;
    setQuestion(value: string): void;

    getType(): PollSchema.QuestionTypeMap[keyof PollSchema.QuestionTypeMap];
    setType(value: PollSchema.QuestionTypeMap[keyof PollSchema.QuestionTypeMap]): void;

    getAnswer(): string;
    setAnswer(value: string): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): QA.AsObject;
    static toObject(includeInstance: boolean, msg: QA): QA.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: QA, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): QA;
    static deserializeBinaryFromReader(message: QA, reader: jspb.BinaryReader): QA;
  }

  export namespace QA {
    export type AsObject = {
      question: string,
      type: PollSchema.QuestionTypeMap[keyof PollSchema.QuestionTypeMap],
      answer: string,
    }
  }

  export interface QuestionTypeMap {
    OPEN: 0;
    CHECKBOX: 1;
    CLOSE: 2;
  }

  export const QuestionType: QuestionTypeMap;
}

export class PollQuestion extends jspb.Message {
  getId(): number;
  setId(value: number): void;

  hasSchema(): boolean;
  clearSchema(): void;
  getSchema(): PollSchema | undefined;
  setSchema(value?: PollSchema): void;

  clearTokensList(): void;
  getTokensList(): Array<VoteToken>;
  setTokensList(value: Array<VoteToken>): void;
  addTokens(value?: VoteToken, index?: number): VoteToken;

  clearVotesList(): void;
  getVotesList(): Array<PollAnswer>;
  setVotesList(value: Array<PollAnswer>): void;
  addVotes(value?: PollAnswer, index?: number): PollAnswer;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PollQuestion.AsObject;
  static toObject(includeInstance: boolean, msg: PollQuestion): PollQuestion.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: PollQuestion, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PollQuestion;
  static deserializeBinaryFromReader(message: PollQuestion, reader: jspb.BinaryReader): PollQuestion;
}

export namespace PollQuestion {
  export type AsObject = {
    id: number,
    schema?: PollSchema.AsObject,
    tokensList: Array<VoteToken.AsObject>,
    votesList: Array<PollAnswer.AsObject>,
  }
}

export class PublicKey extends jspb.Message {
  getKey(): Uint8Array | string;
  getKey_asU8(): Uint8Array;
  getKey_asB64(): string;
  setKey(value: Uint8Array | string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PublicKey.AsObject;
  static toObject(includeInstance: boolean, msg: PublicKey): PublicKey.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: PublicKey, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PublicKey;
  static deserializeBinaryFromReader(message: PublicKey, reader: jspb.BinaryReader): PublicKey;
}

export namespace PublicKey {
  export type AsObject = {
    key: Uint8Array | string,
  }
}

export class RSASignature extends jspb.Message {
  getBallot(): Uint8Array | string;
  getBallot_asU8(): Uint8Array;
  getBallot_asB64(): string;
  setBallot(value: Uint8Array | string): void;

  getSign(): Uint8Array | string;
  getSign_asU8(): Uint8Array;
  getSign_asB64(): string;
  setSign(value: Uint8Array | string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RSASignature.AsObject;
  static toObject(includeInstance: boolean, msg: RSASignature): RSASignature.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RSASignature, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RSASignature;
  static deserializeBinaryFromReader(message: RSASignature, reader: jspb.BinaryReader): RSASignature;
}

export namespace RSASignature {
  export type AsObject = {
    ballot: Uint8Array | string,
    sign: Uint8Array | string,
  }
}

export class SignedEnvelope extends jspb.Message {
  getEnvelope(): Uint8Array | string;
  getEnvelope_asU8(): Uint8Array;
  getEnvelope_asB64(): string;
  setEnvelope(value: Uint8Array | string): void;

  getSign(): Uint8Array | string;
  getSign_asU8(): Uint8Array;
  getSign_asB64(): string;
  setSign(value: Uint8Array | string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SignedEnvelope.AsObject;
  static toObject(includeInstance: boolean, msg: SignedEnvelope): SignedEnvelope.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: SignedEnvelope, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SignedEnvelope;
  static deserializeBinaryFromReader(message: SignedEnvelope, reader: jspb.BinaryReader): SignedEnvelope;
}

export namespace SignedEnvelope {
  export type AsObject = {
    envelope: Uint8Array | string,
    sign: Uint8Array | string,
  }
}

export class VoteReply extends jspb.Message {
  getMess(): string;
  setMess(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VoteReply.AsObject;
  static toObject(includeInstance: boolean, msg: VoteReply): VoteReply.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: VoteReply, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VoteReply;
  static deserializeBinaryFromReader(message: VoteReply, reader: jspb.BinaryReader): VoteReply;
}

export namespace VoteReply {
  export type AsObject = {
    mess: string,
  }
}

export class VoteRequest extends jspb.Message {
  getPollid(): number;
  setPollid(value: number): void;

  hasAnswers(): boolean;
  clearAnswers(): void;
  getAnswers(): PollSchema | undefined;
  setAnswers(value?: PollSchema): void;

  hasSign(): boolean;
  clearSign(): void;
  getSign(): RSASignature | undefined;
  setSign(value?: RSASignature): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VoteRequest.AsObject;
  static toObject(includeInstance: boolean, msg: VoteRequest): VoteRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: VoteRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VoteRequest;
  static deserializeBinaryFromReader(message: VoteRequest, reader: jspb.BinaryReader): VoteRequest;
}

export namespace VoteRequest {
  export type AsObject = {
    pollid: number,
    answers?: PollSchema.AsObject,
    sign?: RSASignature.AsObject,
  }
}

export class VoteToken extends jspb.Message {
  getToken(): Uint8Array | string;
  getToken_asU8(): Uint8Array;
  getToken_asB64(): string;
  setToken(value: Uint8Array | string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VoteToken.AsObject;
  static toObject(includeInstance: boolean, msg: VoteToken): VoteToken.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: VoteToken, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VoteToken;
  static deserializeBinaryFromReader(message: VoteToken, reader: jspb.BinaryReader): VoteToken;
}

export namespace VoteToken {
  export type AsObject = {
    token: Uint8Array | string,
  }
}

