//nolint:lll,gosec
package wordofwisdom

import (
	"math/rand"
	"time"
)

const (
	dictionarySize = 55
)

var dictionary = map[int]string{
	0:  "Wherefore, be not weary in well-doing, for ye are laying the foundation of a great work. And out of small things proceedeth that which is great. - Doctrine and Covenants 64:33",
	1:  "I am Alpha and Omega, Christ the Lord; yea, even I am he, the beginning and the end, the Redeemer of the world. - Doctrine and Covenants 19:1",
	2:  "For behold, it is not meet that I should command in all things; for he that is compelled in all things, the same is a slothful and not a wise servant; wherefore he receiveth no reward. - Doctrine and Covenants 58:26",
	3:  "If thou art merry, praise the Lord with singing, with music, with dancing, and with a prayer of praise and thanksgiving. - Doctrine and Covenants 136:28",
	4:  "And if ye are faithful, ye shall overcome all things, and shall be lifted up at the last day. - Doctrine and Covenants 20:17",
	5:  "Organize yourselves; prepare every needful thing, and establish a house, even a house of prayer, a house of fasting, a house of faith, a house of learning, a house of glory, a house of order, a house of God. - Doctrine and Covenants 88:119",
	6:  "For he will give unto the faithful line upon line, precept upon precept; and I will try you and prove you herewith. - Doctrine and Covenants 98:12",
	7:  "That which is of God is light; and he that receiveth light, and continueth in God, receiveth more light; and that light groweth brighter and brighter until the perfect day. - Doctrine and Covenants 50:24",
	8:  "And faith, hope, charity and love, with an eye single to the glory of God, qualify him for the work. - Doctrine and Covenants 4:5",
	9:  "For he will give unto the faithful line upon line, precept upon precept; and I will try you and prove you herewith. - Doctrine and Covenants 98:12",
	10: "And I give unto you a commandment that you shall teach one another the doctrine of the kingdom. - Doctrine and Covenants 88:77",
	11: "For intelligence cleaveth unto intelligence; wisdom receiveth wisdom; truth embraceth truth; virtue loveth virtue; light cleaveth unto light; mercy hath compassion on mercy and claimeth her own; justice continueth its course and claimeth its own; judgment goeth before the face of him who sitteth upon the throne and governeth and executeth all things. - Doctrine and Covenants 88:40",
	12: "And if ye seek the riches which it is the will of the Father to give unto you, ye shall be the richest of all people, for ye shall have the riches of eternity; and it must needs be that the riches of the earth are mine to give; but beware of pride, lest ye become as the Nephites of old. - Doctrine and Covenants 38:39",
	13: "And see that there is no iniquity in the church, neither hardness with each other, neither lying, backbiting, nor evil speaking; - Doctrine and Covenants 20:54",
	14: "And it shall come to pass that those that die in me shall not taste of death, for it shall be sweet unto them; - Doctrine and Covenants 42:46",
	15: "Therefore, sanctify yourselves that your minds become single to God, and the days will come that you shall see him; for he will unveil his face unto you, and it shall be in his own time, and in his own way, and according to his own will. - Doctrine and Covenants 88:68",
	16: "And if a person gains more knowledge and intelligence in this life through his diligence and obedience than another, he will have so much the advantage in the world to come. - Doctrine and Covenants 130:19",
	17: "And again, I say unto you, I give unto you a commandment, that every man, both elder, priest, teacher, and also member, go to with his might, with the labor of his hands, to prepare and accomplish the things which I have commanded. - Doctrine and Covenants 88:119",
	18: "He that is faithful and endureth shall overcome the world. - Doctrine and Covenants 63:47",
	19: "That which doth not edify is not of God, and is darkness. - Doctrine and Covenants 50:23",
	20: "And this is the gospel, the glad tidings, which the voice out of the heavens bore record unto us - Doctrine and Covenants 76:40",
	21: "And inasmuch as they sought wisdom they might be instructed - Doctrine and Covenants 101:32",
	22: "And that which doth not edify is not of God, and is darkness - Doctrine and Covenants 50:23",
	23: "And that which doth not edify is not of God, and is darkness - Doctrine and Covenants 50:23",
	24: "And that which doth not edify is not of God, and is darkness - Doctrine and Covenants 50:23",
	25: "And if your eye be single to my glory, your whole bodies shall be filled with light, and there shall be no darkness in you; and that body which is filled with light comprehendeth all things - Doctrine and Covenants 88:67",
	26: "And there are many yet on the earth among all sects, parties, and denominations, who are blinded by the subtle craftiness of men, whereby they lie in wait to deceive, and who are only kept from the truth because they know not where to find it - Doctrine and Covenants 123:12",
	27: "And whatsoever ye shall ask the Father in my name, which is right, believing that ye shall receive, behold it shall be given unto you - Doctrine and Covenants 18:20",
	28: "Behold, this is your work, to keep my commandments, yea, with all your might, mind and strength - Doctrine and Covenants 11:20",
	29: "But I have commanded you to bring up your children in light and truth - Doctrine and Covenants 93:40",
	30: "But I have commanded you to bring up your children in light and truth - Doctrine and Covenants 93:40",
	31: "But the laborer in Zion shall labor for Zion; for if they labor for money they shall perish - Doctrine and Covenants 64:39",
	32: "But they are to be used with judgment, not to excess - Doctrine and Covenants 89:7",
	33: "But with some I am not well pleased, for they will not open their mouths, but they hide the talent which I have given unto them, because of the fear of man. Wo unto such, for mine anger is kindled against them - Doctrine and Covenants 60:2",
	34: "Cease to contend one with another; cease to speak evil one of another - Doctrine and Covenants 136:23",
	35: "Come unto me, O ye house of Israel, and it shall be made manifest unto you how great things the Father hath laid up for you, from the foundation of the world; and it hath not come unto you, because of unbelief. - Doctrine and Covenants 121:26",
	36: "Counsel with the Lord in all thy doings, and he will direct thee for good; yea, when thou liest down at night lie down unto the Lord, that he may watch over you in your sleep; and when thou risest in the morning let thy heart be full of thanks unto God; and if ye do these things, ye shall be lifted up at the last day. - Doctrine and Covenants 124:42",
	37: "Cursed are all those that shall lift up the heel against mine anointed, saith the Lord, and cry they have sinned when they have not sinned before me, saith the Lord, but have done that which was meet in mine eyes, and which I commanded them - Doctrine and Covenants 121:16",
	38: "Cursed is he that putteth his trust in man, or maketh flesh his arm, or shall hearken unto the precepts of men, save their precepts shall be given by the power of the Holy Ghost - Doctrine and Covenants 1:19",
	39: "Despair cometh because of iniquity - Doctrine and Covenants 18:15",
	40: "Do not run faster or labor more than you have strength and means provided to enable you to translate; but be diligent unto the end. - Doctrine and Covenants 10:4",
	41: "For behold, this is my work and my glory - Doctrine and Covenants 14:7",
	42: "For behold, this is my work and my glory - Doctrine and Covenants 14:7",
	43: "For behold, this is my work and my glory - Doctrine and Covenants 14:7",
	44: "For behold, this is my work and my glory - Doctrine and Covenants 14:7",
	45: "For the earth is full, and there is enough and to spare; yea, I prepared all things, and have given unto the children of men to be agents unto themselves. - Doctrine and Covenants 104:17",
	46: "For the power is in them, wherein they are agents unto themselves. - Doctrine and Covenants 58:28",
	47: "For those that live shall inherit the earth, and those that die shall rest from all their labors, and their works shall follow them; and they shall receive a crown in the mansions of my Father, which I have prepared for them. - Doctrine and Covenants 59:2",
	48: "For thus shall my church be called in the last days, even The Church of Jesus Christ of Latter-day Saints. - Doctrine and Covenants 115:4",
	49: "For verily I say unto you, that great things await you - Doctrine and Covenants 45:62",
	50: "For verily I say unto you, that great things await you. - Doctrine and Covenants 45:62",
	51: "For verily, verily I say unto you, he that hath the spirit of contention is not of me, but is of the devil, who is the father of contention, and he stirreth up the hearts of men to contend with anger, one with another. - Doctrine and Covenants 76:33",
	52: "For ye are lawful heirs, according to the flesh, and have been hid from the world with Christ in God - Doctrine and Covenants 86:9",
	53: "For ye shall live by every word that proceedeth forth out of the mouth of God. - Doctrine and Covenants 84:44",
	54: "For you, for you will I do it, that others may be benefited - Doctrine and Covenants 11:25",
	55: "For, behold, I say unto you that it mattereth not what ye shall eat or what ye shall drink when ye partake of the sacrament, if it so be that ye do it with an eye single to my glory - Doctrine and Covenants 27:2",
}

func GetRandomQuote() string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	index := rand.Intn(dictionarySize)

	return dictionary[index]
}
