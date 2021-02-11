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
import { host } from '../host';
import { pki, md } from 'node-forge';
import * as bigInt from 'big-integer';

@Component({
  selector: 'vote',
  templateUrl: './vote.component.html',
})
export class VoteComponent {
  questionsList: PollSchema.QA.AsObject[];
  publickey;

  token: string;

  ballot;//BigInteger
  r;//BigInteger

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
        console.log("pollInit.onEnd.status", status, statusMessage);
        console.log("pollInit.onEnd.headers", headers);
        if (status === grpc.Code.OK && message) {
          console.log("pollInit.onEnd.message", message.toObject());
          let pollwithkey: PollWithPublicKey.AsObject = (<PollWithPublicKey>message).toObject()
          this.questionsList = pollwithkey.poll.questionsList;
          this.publickey = pki.publicKeyFromPem("-----BEGIN RSA PUBLIC KEY-----\n" + 
                                         pollwithkey.key.key.toString() + 
                                         "-----END RSA PUBLIC KEY-----");
        }
        console.log("pollInit.onEnd.trailers", trailers);
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

  get diagnostic() { return JSON.stringify(this.questionsList); }

  onSubmit() {
    if (confirm('Czy chcesz wysłać odpowiedź?')) {
      for(let qa of this.questionsList){
        if(qa.type==PollSchema.QuestionType.CLOSE){
          qa.answersList = qa.answersList.map((ans, index) => {
            return parseInt(qa.answersList[0])==index?"true":"false";
          })
        }
        if(qa.type==PollSchema.QuestionType.CHECKBOX){
          qa.answersList = qa.answersList.map((ans, index) => {
            return ans==""?"0":ans;
          })
        }
      }
      console.log(JSON.stringify(this.questionsList));
      let envelope = this.calculateEnvelope()
      let request: EnvelopeToSign = new EnvelopeToSign();
    
      request.setEnvelope(hexToBase64(envelope.toString(16)))
      request.setPollid(this.pollid)
      request.setToken(this.token)
      grpc.unary(Query.SignBallot, {
        request: request,
        host: host,
        onEnd: res => {
          const { status, statusMessage, headers, message, trailers } = res;
          console.log("pollInit.onEnd.status", status, statusMessage);
          console.log("pollInit.onEnd.headers", headers);
          if (status === grpc.Code.OK && message) {
            console.log("pollInit.onEnd.message", message.toObject());
            let signedEnvelope: SignedEnvelope.AsObject = (<SignedEnvelope>message).toObject()
            console.log("smib64: ", signedEnvelope)
            this.sendVote(signedEnvelope)
          }
          else {
            alert("Błąd\n" + statusMessage)
          }
          console.log("pollInit.onEnd.trailers", trailers);
        }
      });
    }
  }

  calculateEnvelope(){
    // Generate ballot to be signed.
    let N = bigInt(this.publickey.n.toString())
    let e = bigInt(this.publickey.e.toString()) // Should be always 65537
    this.ballot = bigInt.randBetween(bigInt("2"), N);

    // We are hashing ballot.
    let sha256 = md.sha256.create();
    sha256.update(this.ballot.toString())
    let hash = sha256.digest().toHex()
    let m = bigInt(hash, 16)

    // Get random blinding factor.
    this.r = bigInt.randBetween(bigInt("2"), N);

    // We want to send m*r^e mod N to server.
    let re = this.r.modPow(e, N)

    // blinded = m*(r^e) mod N
    let envelope = (m.multiply(re)).mod(N)
    return envelope
  }

  sendVote(senv: SignedEnvelope.AsObject) {
    // Having (m^d)*r mod N we are removing blinding factor r,
    let N = bigInt(this.publickey.n.toString())
    let hex = base64ToHex(senv.sign)
    let smi = bigInt(hex, 16)
    let revr = this.r.modInv(N)
    let smirevr = smi.multiply(revr)
    // Now we can calculate second part of sign.
    // sign = smirevr mod N = m^d mod N
    let sign = smirevr.mod(N)
    
    let request: VoteRequest = new VoteRequest();
    let signature: RSASignature = new RSASignature();
    let schema: PollSchema = this.toSchema();
    
    request.setPollid(this.pollid)
    signature.setBallot(hexToBase64(this.ballot.toString(16)))
    signature.setSign(hexToBase64(sign.toString(16)))
    request.setAnswers(schema)
    request.setSign(signature)
    grpc.unary(Query.PollVote, {
      request: request,
      host: host,
      onEnd: res => {
        const { status, statusMessage, headers, message, trailers } = res;
        console.log("pollInit.onEnd.status", status, statusMessage);
        console.log("pollInit.onEnd.headers", headers);
        if (status === grpc.Code.OK && message) {
          console.log("pollInit.onEnd.message", message.toObject());
          this.router.navigate(['/results', this.pollid]);
        }
        console.log("pollInit.onEnd.trailers", trailers);
      }
    });
  }

  toSchema() {
    let schema = new PollSchema();
    for (let qa of this.questionsList){
      const QA = new PollSchema.QA();
      QA.setQuestion(qa.question);
      QA.setType(qa.type as 0 | 1 | 2);
      QA.setOptionsList(qa.optionsList);
      QA.setAnswersList(qa.answersList);
      schema.addQuestions(QA);
    }
    return schema
  }
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

function hexToBase64(hexstring) {
    return btoa(hexstring.match(/\w{2}/g).map(function(a) {
        return String.fromCharCode(parseInt(a, 16));
    }).join(""));
}
