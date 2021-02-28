import { Component } from '@angular/core';
import { ActivatedRoute,
         Router } from '@angular/router';

import { grpc } from '@improbable-eng/grpc-web';
import { Query } from "Projekt_Rada/query/query_pb_service";
import { EnvelopeToSign, 
         GetPollRequest, 
         PollSchema, 
         PollWithPublicKey, 
         RSASignature,
         SignedEnvelope, 
         VoteRequest } from "Projekt_Rada/query/query_pb";
import { QAListToSchema,
         toEnvelope,
         toRSASignature,
         toVoteRequest } from "../proto_parsing";
import { host } from '../host';
import { pki, md } from 'node-forge';
import * as bigInt from 'big-integer';

@Component({
  selector: 'vote',
  templateUrl: './vote.component.html',
})
export class VoteComponent {
  questionsList: PollSchema.QA.AsObject[];
  
  publickey; //PublicKey

  token: string;

  ballot: bigInt.BigInteger;
  r: bigInt.BigInteger;

  pollid: number;
  inpid: number;

  constructor (private route: ActivatedRoute, private router: Router) {
    this.pollid = parseInt(this.route.snapshot.paramMap.get('pollid'));
    if(!isNaN(this.pollid)){
      let request: GetPollRequest = new GetPollRequest();
      request.setPollid(this.pollid);
      grpc.unary(Query.GetPoll, {
        request: request,
        host: host,
        onEnd: res => {
          const { status, statusMessage, headers, message, trailers } = res;
          if (status === grpc.Code.OK && message) {
            let pollwithkey: PollWithPublicKey.AsObject = (<PollWithPublicKey>message).toObject()
            this.questionsList = pollwithkey.poll.questionsList;
            let pempublickey = "-----BEGIN RSA PUBLIC KEY-----\n" +
                               pollwithkey.key.key.toString() +
                               "-----END RSA PUBLIC KEY-----";
            this.publickey = pki.publicKeyFromPem(pempublickey);
          }
        }
      });
    }
  }

  getPollid() {
    this.router.navigate(['/vote', this.inpid]);
  }

  trackOption(index: number, option: string) {
    return index;
  }

  onSubmit() {
    if (confirm('Czy chcesz wysłać odpowiedź?')) {
      // Parse input
      for(let qa of this.questionsList){
        // If question is close, then first answer field is containing id of selected answer.
        if(qa.type==PollSchema.QuestionType.CLOSE){
          qa.answersList = qa.answersList.map((ans, index) => {
            return parseInt(qa.answersList[0])==index?"true":"false";
          })
        }
        // If question is checkbox, we want to replace empty strings to false answer
        if(qa.type==PollSchema.QuestionType.CHECKBOX){
          qa.answersList = qa.answersList.map((ans, index) => {
            return ans==""?"false":ans;
          })
        }
      }

      let envelope = this.calculateEnvelope()
      let request: EnvelopeToSign = toEnvelope(envelope.toString(16), this.pollid, this.token);
    
      grpc.unary(Query.SignBallot, {
        request: request,
        host: host,
        onEnd: res => {
          const { status, statusMessage, headers, message, trailers } = res;
          if (status === grpc.Code.OK && message) {
            let signedEnvelope: SignedEnvelope.AsObject = (<SignedEnvelope>message).toObject()
            this.sendVote(signedEnvelope)
          }
          else {
            alert("Błąd\n" + statusMessage)
          }
        }
      });
    }
  }

  sendVote(senv: SignedEnvelope.AsObject) {
    let sign = this.calculateSign(senv.sign);

    let schema: PollSchema = QAListToSchema(this.questionsList);
    let signature: RSASignature = toRSASignature(this.ballot.toString(16), sign.toString(16));
    let request: VoteRequest = toVoteRequest(this.pollid, schema, signature);
    
    grpc.unary(Query.PollVote, {
      request: request,
      host: host,
      onEnd: res => {
        const { status, statusMessage, headers, message, trailers } = res;
        if (status === grpc.Code.OK && message) {
          this.router.navigate(['/results', this.pollid]);
        }
      }
    });
  }

  calculateEnvelope(){
    // Generate ballot to be signed.
    let N = bigInt.call({}, this.publickey.n.toString())
    let e = bigInt.call({}, this.publickey.e.toString()) // Should be always 65537
    this.ballot = bigInt.randBetween(bigInt.call({}, "2"), N);

    // We are hashing ballot.
    let sha256 = md.sha256.create();
    sha256.update(this.ballot.toString())
    let hash = sha256.digest().toHex()
    let m = bigInt.call({}, hash, 16)

    // Get random blinding factor.
    this.r = bigInt.randBetween(bigInt.call({}, "2"), N);

    // We want to send m*r^e mod N to server.
    let re = this.r.modPow(e, N)

    // blinded = m*(r^e) mod N
    let envelope = (m.multiply(re)).mod(N)
    return envelope
  }

  calculateSign(signedEnvelope: string | Uint8Array) {
    // Having (m^d)*r mod N we are removing blinding factor r,
    let N = bigInt.call({}, this.publickey.n.toString())
    let sm = bigInt.call({}, base64ToHex(signedEnvelope), 16)
    let revr = this.r.modInv(N)
    let smrevr = sm.multiply(revr)
    // Now we can calculate second part of sign.
    // sign = smrevr mod N = m^d mod N
    let sign = smrevr.mod(N)
    return sign;
  }

  get diagnostic() { return JSON.stringify(this.questionsList); }
}

function base64ToHex(str) {
  const raw = atob(str);
  let result = '';
  for (let i = 0; i < raw.length; i++) {
    const hex = raw.charCodeAt(i).toString(16);
    result += (hex.length === 2 ? hex : '0' + hex);
  }
  return result.toLowerCase();
}
