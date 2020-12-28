const Mfrc522 = require('mfrc522-rpi')
const SoftSPI = require('rpi-softspi')
const logger = require('../helper/logger')(module)
const audioBookPlayerService = require('../service/audio-book-player-service')

let mfrc522

const maxHits = 3
let hitCounter = 0
let lastUid = null

const initRfid = () => {
    logger.info('initializing card reader...')
    mfrc522 = initMfrc522()

    logger.info('starting scan interval..')
    setInterval(scanForCards, 500)
}

const initMfrc522 = () => {
    const softSPI = new SoftSPI({
        clock: 23, // pin number of SCLK
        mosi: 19, // pin number of MOSI
        miso: 21, // pin number of MISO
        client: 24 // pin number of CS
    })

    // GPIO 24 can be used for buzzer bin (PIN 18), Reset pin is (PIN 22).
    // I believe that channing pattern is better for configuring pins which are optional methods to use.
    return new Mfrc522(softSPI).setResetPin(22)
}

const scanForCards = () => {
    // reset card
    mfrc522.reset()

    // Scan for cards
    let response = mfrc522.findCard()
    if (!response.status) {
        checkHits(null)
        return
    }
    // logger.info('Card detected, CardType: ' + response.bitSize)

    // Get the UID of the card
    response = mfrc522.getUid()
    if (!response.status) {
        logger.info('UID Scan Error')
        checkHits(null)
        return
    }
    // If we have the UID, continue
    const uid = response.data
    // logger.info('Card uid: ' + uidToString(uid))

    checkHits(uid)
}

const uidToString = (uid) => {
    if (!uid) {
        return null
    }

    return uid.reduce((s, b) => {
        return s + (b < 16 ? '0' : '') + b.toString(16)
    }, '')
}

const checkHits = (sourceUid) => {
    const uid = uidToString(sourceUid)

    // logger.debug('got uid: ' + uid)
    if (lastUid !== uid) {
        // logger.debug('change detected: From ' + lastUid + ' -> type: ' + typeof lastUid + ' to ' + uid + ' -> type: ' + typeof uid)
        if (hitCounter >= maxHits) {
            // logger.debug('reached max hitCount')
            lastUid = uid || null
            hitCounter = 0

            if (!uid) {
                logger.debug('No card detected fired')
                audioBookPlayerService.nfcNoCardDetected()
            } else {
                logger.debug('Card detected fired')
                audioBookPlayerService.nfcCardDetected(uid)
            }
        } else {
            hitCounter++
            // logger.debug('Max hitCount not reached. HitCounter increased to: ' + hitCounter)
        }
    } else {
        // logger.debug('no change detected')
        hitCounter = 0
    }
}

module.exports = {
    initRfid
}
