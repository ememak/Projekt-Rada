// package: query
// file: query/query.proto

import * as query_query_pb from "../query/query_pb";
import {grpc} from "@improbable-eng/grpc-web";

type QueryGetPoll = {
  readonly methodName: string;
  readonly service: typeof Query;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof query_query_pb.GetPollRequest;
  readonly responseType: typeof query_query_pb.PollWithPublicKey;
};

type QueryPollInit = {
  readonly methodName: string;
  readonly service: typeof Query;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof query_query_pb.PollSchema;
  readonly responseType: typeof query_query_pb.PollQuestion;
};

type QuerySignBallot = {
  readonly methodName: string;
  readonly service: typeof Query;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof query_query_pb.EnvelopeToSign;
  readonly responseType: typeof query_query_pb.SignedEnvelope;
};

type QueryPollVote = {
  readonly methodName: string;
  readonly service: typeof Query;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof query_query_pb.VoteRequest;
  readonly responseType: typeof query_query_pb.VoteReply;
};

export class Query {
  static readonly serviceName: string;
  static readonly GetPoll: QueryGetPoll;
  static readonly PollInit: QueryPollInit;
  static readonly SignBallot: QuerySignBallot;
  static readonly PollVote: QueryPollVote;
}

export type ServiceError = { message: string, code: number; metadata: grpc.Metadata }
export type Status = { details: string, code: number; metadata: grpc.Metadata }

interface UnaryResponse {
  cancel(): void;
}
interface ResponseStream<T> {
  cancel(): void;
  on(type: 'data', handler: (message: T) => void): ResponseStream<T>;
  on(type: 'end', handler: (status?: Status) => void): ResponseStream<T>;
  on(type: 'status', handler: (status: Status) => void): ResponseStream<T>;
}
interface RequestStream<T> {
  write(message: T): RequestStream<T>;
  end(): void;
  cancel(): void;
  on(type: 'end', handler: (status?: Status) => void): RequestStream<T>;
  on(type: 'status', handler: (status: Status) => void): RequestStream<T>;
}
interface BidirectionalStream<ReqT, ResT> {
  write(message: ReqT): BidirectionalStream<ReqT, ResT>;
  end(): void;
  cancel(): void;
  on(type: 'data', handler: (message: ResT) => void): BidirectionalStream<ReqT, ResT>;
  on(type: 'end', handler: (status?: Status) => void): BidirectionalStream<ReqT, ResT>;
  on(type: 'status', handler: (status: Status) => void): BidirectionalStream<ReqT, ResT>;
}

export class QueryClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  getPoll(
    requestMessage: query_query_pb.GetPollRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: query_query_pb.PollWithPublicKey|null) => void
  ): UnaryResponse;
  getPoll(
    requestMessage: query_query_pb.GetPollRequest,
    callback: (error: ServiceError|null, responseMessage: query_query_pb.PollWithPublicKey|null) => void
  ): UnaryResponse;
  pollInit(
    requestMessage: query_query_pb.PollSchema,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: query_query_pb.PollQuestion|null) => void
  ): UnaryResponse;
  pollInit(
    requestMessage: query_query_pb.PollSchema,
    callback: (error: ServiceError|null, responseMessage: query_query_pb.PollQuestion|null) => void
  ): UnaryResponse;
  signBallot(
    requestMessage: query_query_pb.EnvelopeToSign,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: query_query_pb.SignedEnvelope|null) => void
  ): UnaryResponse;
  signBallot(
    requestMessage: query_query_pb.EnvelopeToSign,
    callback: (error: ServiceError|null, responseMessage: query_query_pb.SignedEnvelope|null) => void
  ): UnaryResponse;
  pollVote(
    requestMessage: query_query_pb.VoteRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: query_query_pb.VoteReply|null) => void
  ): UnaryResponse;
  pollVote(
    requestMessage: query_query_pb.VoteRequest,
    callback: (error: ServiceError|null, responseMessage: query_query_pb.VoteReply|null) => void
  ): UnaryResponse;
}

