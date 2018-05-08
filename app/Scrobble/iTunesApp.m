//
//  iTunes.m
//  Scrobble
//
//  Created by Nishanth Shanmugham on 5/8/18.
//  Copyright © 2018 Nishanth Shanmugham. All rights reserved.
//

#import <Foundation/Foundation.h>
#import "iTunesApp.h"

@implementation iTunesApp : NSObject

- (id) init {
    self = [super init];
    if (self) {
        _app = [SBApplication applicationWithBundleIdentifier:@"com.apple.iTunes"];
    }
    return self;
}

- (BOOL) currentTrackIsSong {
    return [[[self app] currentTrack] mediaKind] == iTunesEMdKSong;
}

@end
